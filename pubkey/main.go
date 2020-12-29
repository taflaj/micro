// pubkey: main.go

package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/taflaj/micro/messaging"
	"github.com/taflaj/micro/pubkey/models"
)

const (
	name    = "pubkey"
	port    = "8003"
	version = "0.1.0"
)

// Env contains the database access environment
type Env struct {
	db models.DataStore
}

var env *Env

func init() {
	log.SetFlags(log.Flags() | log.Lmicroseconds)
}

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func logIt(r *http.Request) {
	log.Printf("%v %v from %v using %v", r.Method, r.URL.Path, r.RemoteAddr, r.Header["User-Agent"][0])
}

func fail(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	fmt.Fprint(w, err)
}

func handler(w http.ResponseWriter, r *http.Request) {
	logIt(r)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	var err error
	msg, err := messaging.GetMessage(r)
	if err != nil {
		fail(w, http.StatusInternalServerError, err)
	} else if len(msg.Command) != 3 {
		fail(w, http.StatusBadRequest, errors.New("Bad request format"))
	}
	var result string
	key := msg.Command[2]
	if key == "version" {
		result = version
	} else {
		result, err = env.db.GetPublicKey(key)
		// log.Printf("%v; %v; %v", key, result, err)
		if err != nil {
			log.Printf("%v", err)
			fail(w, http.StatusBadRequest, fmt.Errorf("%v not found", key))
		}
	}
	fmt.Fprint(w, result)
}

func run(file string) {
	db, err := models.Open(file)
	check(err)
	env = &Env{db}
	defer env.db.Close()
	http.HandleFunc("/", handler)
	go func() {
		_, err := messaging.Get(messaging.Messenger, "register/"+name+"/"+port)
		check(err)
	}()
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please include database file.")
	} else {
		run(os.Args[1])
	}
}
