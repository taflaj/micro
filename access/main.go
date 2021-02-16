// access: main.go

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/taflaj/micro/access/models"
	"github.com/taflaj/micro/modules/logger"
	"github.com/taflaj/micro/modules/messaging"
)

const (
	name    = messaging.Access
	port    = "9995"
	version = "0.1.0 dev"
)

// Env contains the database access environment
type Env struct {
	db models.DataStore
}

var env *Env

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
	w.Header().Set("Content-Type", "text/plain")
	var err error
	msg, err := messaging.GetMessage(r)
	if err != nil {
		// messaging.Fail(w, http.StatusInternalServerError, err)
		messaging.FailInternal(w, err)
	} else if msg.IP == -1 {
		// logger.GetLogger().Printf("DEBUG msg.Data=%#v", msg.Data)
		data := []byte(msg.Data.(string))
		query := messaging.AccessLevelQuery{}
		json.Unmarshal(data, &query)
		logger.GetLogger().Printf("DEBUG query=%#v", query)
		level, _ := env.db.GetAccess(query.Address, query.Service)
		logger.GetLogger().Printf("DEBUG level=%#v", level)
		// query := messaging.AccessLevelQuery{}
		// 	Address: data["Address"].(int),
		// 	Service: data["Service"].(string),
		// }
		// err = json.Unmarshal([]byte(msg.Data.([]byte)), &query)
		// logger.GetLogger().Printf("DEBUG query=%#v, %v", query, err)
		// w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(level)
		// messaging.Ok(w)
		fmt.Fprintf(w, "%s", response)
	} else if len(msg.Command) != 3 {
		messaging.FailBadRequestStandard(w)
	} else {
		// var result string
		service := msg.Command[2]
		if service == "version" {
			// result = version
			fmt.Fprint(w, version)
		} else if msg.Command[0] != "set" {
			messaging.FailBadRequestStandard(w)
		} else {
			switch msg.Method {
			case http.MethodDelete:
				// resetAccess(w, msg)
				env.db.ResetAccess(w, msg)
			case http.MethodPut:
				env.db.SetAccess(w, msg)
				// if err = env.db.SetAccess(w, msg); err != nil {
				// 	messaging.Fail(w, http.StatusInternalServerError, err)
				// }
			default:
				messaging.FailBadRequestStandard(w)
			}
			// result = "Ok"
		}
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
