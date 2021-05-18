// access: main.go

package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/taflaj/services/access/models"
	"github.com/taflaj/services/modules/logger"
	"github.com/taflaj/services/modules/messaging"
	"github.com/taflaj/util/server"
)

const (
	name    = messaging.Access
	port    = "9995"
	version = "0.1.1"
)

// Env contains the database access environment
type Env struct {
	db models.DataStore
}

var env *Env

func init() {
	logger.NewLogger(name, logger.Debug)
}

func logIt(r *http.Request) {
	who := r.Header["User-Agent"][0]
	logger.GetLogger().Logdf(logger.Debug, 1, "%#v", r)
	logger.GetLogger().Logdf(logger.Info, 1, "%v %v from %v using %v", r.Method, r.URL.Path, r.RemoteAddr, who)
}

func handler(w http.ResponseWriter, r *http.Request) {
	logIt(r)
	message, err := messaging.GetMessage(r)
	if err != nil {
		messaging.FailInternal(w, err)
	} else if len(message.Command) != 3 {
		messaging.FailBadRequest(w)
	} else {
		command := message.Command[0]
		service := message.Command[2]
		if service == "version" {
			fmt.Fprint(w, version)
		} else {
			level := env.db.GetAccess(message.IP, name)
			data := message.GetData()
			host := (*data)["ip"]
			ip, err := messaging.IPtoInt(host)
			done := make(chan bool, 1)
			if err != nil {
				messaging.FailBadRequest(w, err)
			} else if command == "get" && message.Method == http.MethodGet {
				go func() {
					if level.Defined && level.CanRead {
						logger.GetLogger().Logf(logger.Debug, "%#v", message)
						level = env.db.GetAccess(ip, service)
						if level.Defined {
							fmt.Fprint(w, level.Level)
						} else {
							fmt.Fprint(w, "default")
						}
					} else {
						messaging.FailNotAuthorized(w)
					}
					done <- true
				}()
				<-done
			} else if command == "set" {
				if level.Defined && level.CanWrite {
					go func() {
						if message.Method == http.MethodDelete {
							logger.GetLogger().Logf(logger.Info, "Resetting access for '%v' to '%v'", host, service)
							if err = env.db.ResetAccess(ip, service); err != nil {
								messaging.FailInternal(w, err)
							} else {
								messaging.Ok(w)
							}
						} else if message.Method == http.MethodPut {
							access := (*data)["access"]
							owner := (*data)["owner"]
							remarks, ok := (*data)["remarks"]
							if !ok {
								remarks = fmt.Sprint(time.Now())
							}
							level := messaging.AccessLevel{IP: ip, Service: service, Level: access, Defined: true}
							switch access {
							case "no":
								level.CanRead, level.CanWrite = false, false
							case "ro":
								level.CanRead, level.CanWrite = true, false
							case "rw":
								level.CanRead, level.CanWrite = true, true
							case "wo":
								level.CanRead, level.CanWrite = false, true
							default:
								level.Defined = false
							}
							if level.Defined {
								logger.GetLogger().Logf(logger.Info, "Giving %v (%v) %v access to %v", host, owner, access, service)
								if err = env.db.SetAccess(&level, owner, remarks); err != nil {
									messaging.FailInternal(w, err)
								} else {
									messaging.Ok(w)
								}
							} else {
								messaging.FailBadRequest(w, fmt.Errorf("invalid access \"%v\"", access))
							}
						} else {
							messaging.FailBadRequest(w)
						}
						done <- true
					}()
					<-done
				} else {
					messaging.FailNotAuthorized(w)
				}
			} else {
				messaging.FailBadRequest(w)
			}
		}
	}
}

func run(file string) {
	db, err := models.Open(file)
	if err != nil {
		logger.GetLogger().Log(logger.Critical, err)
		panic(err)
	}
	env = &Env{db}
	defer env.db.Close()
	messaging.Get(messaging.Messenger, "register/"+name+"/"+port+"/localhost")
	var handler = server.Handlers{{Pattern: "/", Handler: handler}}
	me := server.NewServer(":"+port, &handler)
	me.SetOnStart(func() { logger.GetLogger().Logf(logger.Info, "Listening on port %v", port) })
	me.SetOnFail(func(err error) { logger.GetLogger().Log(logger.Error, err) })
	me.SetOnInterrupt(func(signal os.Signal) { logger.GetLogger().Logf(logger.Warning, "Received %v", signal) })
	me.Start()
	logger.GetLogger().Log(logger.Info, "Exiting now")
}

func main() {
	if len(os.Args) < 2 {
		logger.GetLogger().Log(logger.Critical, "Please include database file")
	} else {
		run(os.Args[1])
	}
}
