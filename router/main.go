// router: main.go

package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/taflaj/micro/modules/logger"
	"github.com/taflaj/micro/modules/messaging"
)

type dir map[string]*messaging.Host

var directory dir

const (
	name    = "router"
	port    = "9998"
	version = "0.1.2 dev"
)

var (
	exists   = struct{}{}
	services = map[string]struct{}{
		"access": exists,
		"pubkey": exists,
		"random": exists,
		"router": exists,
		"server": exists,
	}
	log logger.Logger
)

func init() {
	directory = make(dir)
	logger.NewLogger(name)
}

func logIt(r *http.Request) {
	who := r.Header["User-Agent"][0]
	logger.GetLogger().Printf("%v %v from %v using %v", r.Method, r.URL.Path, r.RemoteAddr, who)
	logger.GetLogger().Spy(who)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	// incoming request format: register/<service>/<port>[/<address>]
	logIt(r)
	splits := strings.Split(r.URL.Path, "/")
	if len(splits) > 3 {
		service := splits[2]
		if _, exists := directory[service]; exists {
			w.WriteHeader(http.StatusAlreadyReported)
		} else if _, ok := services[service]; !ok {
			messaging.Fail(w, http.StatusUnauthorized, fmt.Errorf("Service '%v' is not authorized", service))
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
	// logger.GetLogger().Printf("%#v\n  %v", msg, err)
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
					// messaging.Fail(w, http.StatusInternalServerError, err)
					messaging.FailInternal(w, err)
				} else {
					w.WriteHeader(response.Code)
					w.Header().Set("Content-Type", response.ContentType)
					fmt.Fprint(w, response.Payload)
				}
			} else if _, ok := services[service]; ok {
				// messaging.Fail(w, http.StatusInternalServerError, fmt.Errorf("Service '%v' is not active", service))
				messaging.FailInternal(w, fmt.Errorf("Service '%v' is not active", service))
			} else {
				messaging.FailNotAuthorizedStandard(w)
			}
			// fmt.Fprintf(w, "%#v\n", msg)
			// return
		}
	} else {
		messaging.Fail(w, http.StatusInternalServerError, err)
	}
}

func main() {
	http.HandleFunc("/register/", registerHandler)
	http.HandleFunc("/unregister/", unregisterHandler)
	http.HandleFunc("/", generalHandler)
	logger.GetLogger().Printf("Listening on port %v", port)
	logger.GetLogger().Fatal(http.ListenAndServe(":"+port, nil))
}
