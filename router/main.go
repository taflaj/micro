// router: main.go

package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/taflaj/micro/messaging"
)

type dir map[string]*messaging.Host

var directory dir

const (
	name    = "router"
	port    = "8001"
	version = "0.1.1"
)

func init() {
	log.SetFlags(log.Flags() | log.Lmicroseconds)
	directory = make(dir)
}

func logIt(r *http.Request) {
	log.Printf("%v %v from %v using %v", r.Method, r.URL.Path, r.RemoteAddr, r.Header["User-Agent"][0])
}

func fail(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(w, err)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	// incoming request format: register/<service>/<port>
	logIt(r)
	splits := strings.Split(r.URL.Path, "/")
	if len(splits) > 3 {
		service := splits[2]
		if _, exists := directory[service]; exists {
			w.WriteHeader(http.StatusAlreadyReported)
		} else {
			w.WriteHeader(http.StatusOK)
			address := strings.Split(r.RemoteAddr, ":")[0]
			if len(splits) > 4 {
				address = splits[4]
			}
			port := splits[3]
			directory[service] = &messaging.Host{Address: address, Port: port}
		}
		// w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
	fmt.Fprintf(w, "%#v\n", directory)
}

func unregisterHandler(w http.ResponseWriter, r *http.Request) {
	// incoming request format: unregister/<service>
	logIt(r)
	service := strings.Split(r.URL.Path, "/")[2]
	delete(directory, service)
	fmt.Fprintf(w, "%#v\n", directory)
}

func generalHandler(w http.ResponseWriter, r *http.Request) {
	logIt(r)
	var err error
	msg, err := messaging.GetMessage(r)
	// log.Printf("%#v\n  %v", msg, err)
	if err == nil {
		service := msg.Service
		if service == name && msg.Command[0] == "get" && msg.Command[2] == "version" {
			fmt.Fprint(w, version)
		} else {
			host, present := directory[service]
			if present {
				// fmt.Fprintf(w, "Forwarding message to %#v\n", host)
				response, err := msg.Send(host)
				if err != nil {
					fail(w, err)
				} else {
					w.WriteHeader(response.Code)
					w.Header().Set("Content-Type", response.ContentType)
					fmt.Fprint(w, response.Payload)
				}
			} else {
				fail(w, fmt.Errorf("Service '%v' is not active", service))
			}
			// fmt.Fprintf(w, "%#v\n", msg)
			// return
		}
	} else {
		fail(w, err)
	}
}

func main() {
	http.HandleFunc("/register/", registerHandler)
	http.HandleFunc("/unregister/", unregisterHandler)
	http.HandleFunc("/", generalHandler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
