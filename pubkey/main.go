// pubkey: main.go

package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/taflaj/micro/modules/logger"
	"github.com/taflaj/micro/modules/messaging"
	"github.com/taflaj/micro/pubkey/models"
)

const (
	name    = "pubkey"
	port    = "9996"
	version = "0.1.1 dev"
)

// Env contains the database access environment
type Env struct {
	db models.DataStore
}

var (
	env *Env
	// log logger.Logger
)

func init() {
	logger.NewLogger(name)
}

func check(err error) {
	if err != nil {
		logger.GetLogger().Panic(err)
	}
}

func logIt(r *http.Request) {
	who := r.Header["User-Agent"][0]
	logger.GetLogger().Printf("%v %v from %v using %v", r.Method, r.URL.Path, r.RemoteAddr, who)
	// logger.GetLogger().Spy(who)
}

func handler(w http.ResponseWriter, r *http.Request) {
	logIt(r)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	var err error
	msg, err := messaging.GetMessage(r)
	// logger.GetLogger().Printf("%#v", msg.Command)
	if err != nil {
		// messaging.Fail(w, http.StatusInternalServerError, err)
		messaging.FailInternal(w, err)
	} else if len(msg.Command) != 3 {
		messaging.FailBadRequestStandard(w)
	} else {
		// var result string
		key := msg.Command[2]
		command := msg.Command[0]
		if command == "get" && msg.Method == http.MethodGet {
			if key == "version" {
				fmt.Fprint(w, version)
			} else {
				result, err := env.db.GetPublicKey(key)
				// logger.GetLogger().Printf("%v; %v; %v", key, result, err)
				if err == nil {
					fmt.Fprint(w, result)
				} else {
					logger.GetLogger().Printf("%v", err)
					messaging.Fail(w, http.StatusBadRequest, fmt.Errorf("'%v' not found", key))
				}
			}
			// } else if command == "set" {
			// 	switch msg.Method {
			// 	case http.MethodDelete:
			// 		env.db.DeletePublicKey(w, msg, name)
			// 	case http.MethodPut:
			// 		env.db.SetPublicKey(w, msg, name)
			// 	default:
			// 		messaging.FailBadRequestStandard(w)
			// 	}
		} else {
			messaging.FailBadRequestStandard(w)
		}
		// fmt.Fprint(w, result)
	}
}

func run(file string) {
	db, err := models.Open(file)
	check(err)
	env = &Env{db}
	defer env.db.Close()
	http.HandleFunc("/", handler)
	go func() {
		_, err := messaging.Get(messaging.Messenger, "register/"+name+"/"+port+"/localhost")
		check(err)
	}()
	logger.GetLogger().Printf("Listening on port %v", port)
	logger.GetLogger().Fatal(http.ListenAndServe(":"+port, nil))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please include database file.")
	} else {
		run(os.Args[1])
	}
}
