// random: main.go

package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/taflaj/micro/messaging"
	"github.com/taflaj/util/random"
)

const (
	name    = "random"
	port    = "8002"
	version = "0.1.0"
)

func init() {
	log.SetFlags(log.Flags() | log.Lmicroseconds)
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
	msg, err := messaging.GetMessage(r)
	if err != nil {
		fail(w, http.StatusInternalServerError, err)
	} else if len(msg.Command) < 3 {
		fail(w, http.StatusBadRequest, errors.New("Bad request format"))
	} else {
		// log.Printf("%#v", msg)
		length := 32
		if len(msg.Command) > 3 {
			l, err := strconv.Atoi(msg.Command[3])
			if err == nil {
				length = l
			}
		}
		var result string
		t := msg.Command[2]
		switch t {
		case "alpha":
			result, err = random.Alpha(length)
		case "alphanum":
			result, err = random.AlphaNum(length)
		case "any":
			result, err = random.Any(length)
		case "hex":
			result, err = random.Hex(length)
		case "number":
			result, err = random.Number(length)
		case "special":
			result, err = random.Special(length)
		case "version":
			result = version
		default:
			err = fmt.Errorf("Unknown argument \"%v\"", t)
		}
		if err == nil {
			fmt.Fprint(w, result)
		} else {
			fail(w, http.StatusInternalServerError, err)
		}
	}
}

func main() {
	http.HandleFunc("/", handler)
	go func() {
		if _, err := messaging.Get(messaging.Messenger, "register/"+name+"/"+port+"/localhost"); err != nil {
			log.Fatal(err)
		}
	}()
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
