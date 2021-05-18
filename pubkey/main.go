// pubkey: main.go

package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/taflaj/services/modules/logger"
	"github.com/taflaj/services/modules/messaging"
	"github.com/taflaj/services/pubkey/models"
	"github.com/taflaj/util/server"
)

const (
	name    = "pubkey"
	port    = "9996"
	version = "0.1.1"
)

// Env contains the data access environment
type Env struct {
	db models.DataStore
}

var env *Env

func init() {
	logger.NewLogger(name, logger.Info)
}

func logIt(r *http.Request) {
	who := r.Header["User-Agent"][0]
	logger.GetLogger().Logdf(logger.Debug, 1, "%#v", r)
	logger.GetLogger().Logdf(logger.Info, 1, "%v %v from %v using %v", r.Method, r.URL.Path, r.RemoteAddr, who)
}

func handler(w http.ResponseWriter, r *http.Request) {
	logIt(r)
	var err error
	message, err := messaging.GetMessage(r)
	if err != nil {
		messaging.FailInternal(w, err)
	} else if len(message.Command) != 3 {
		messaging.FailBadRequest(w)
	} else {
		key := message.Command[2]
		command := message.Command[0]
		if command == "get" && message.Method == http.MethodGet {
			if key == "version" {
				fmt.Fprint(w, version)
			} else {
				result, err := env.db.GetPublicKey(key)
				if err == nil {
					fmt.Fprint(w, result)
				} else {
					logger.GetLogger().Log(logger.Error, err)
					messaging.FailBadRequest(w, fmt.Errorf("key for '%v' not found", key))
				}
			}
		} else {
			messaging.FailBadRequest(w)
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
