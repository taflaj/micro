// router: main.go

package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/taflaj/services/modules/logger"
	"github.com/taflaj/services/modules/messaging"
	"github.com/taflaj/util/server"
)

const (
	name    = "router"
	port    = "9998"
	version = "0.1.3"
)

type dir map[string]*messaging.Host

var (
	directory dir
	exists    = struct{}{}
	services  = map[string]struct{}{
		"access": exists,
		"digest": exists,
		"pubkey": exists,
		"random": exists,
		"router": exists,
		"server": exists,
	}
)

func init() {
	logger.NewLogger(name, logger.Info)
	// logger.NewLogger(name, logger.Debug)
	directory = make(dir)
}

func (d dir) dump() {
	if logger.GetLogger().GetLevel() <= logger.Debug {
		output := "dir: "
		for k, v := range d {
			output += fmt.Sprintf("{%v@%v:%v}", k, v.Address, v.Port)
		}
		logger.GetLogger().Logd(logger.Debug, 1, output)
	}
}

func logIt(r *http.Request) {
	who := r.Header["User-Agent"][0]
	logger.GetLogger().Logdf(logger.Debug, 1, "%#v", r)
	logger.GetLogger().Logdf(logger.Info, 1, "%v %v from %v using %v", r.Method, r.URL.Path, r.RemoteAddr, who)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	// incoming request format: /register/<service>/<port>[/<address>]
	logIt(r)
	splits := strings.Split(r.URL.Path, "/")
	if len(splits) > 3 {
		service := splits[2]
		if _, exists := directory[service]; exists {
			w.WriteHeader(http.StatusAlreadyReported)
		} else if _, ok := services[service]; !ok {
			logger.GetLogger().Log(logger.Error, fmt.Errorf("service '%s' is not authorized", service))
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			w.WriteHeader(http.StatusOK)
			address := strings.Split(r.RemoteAddr, ":")[0]
			if len(splits) > 4 {
				address = splits[4]
			}
			port := splits[3]
			directory[service] = &messaging.Host{Address: address, Port: port}
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
	directory.dump()
	messaging.Ok(w)
}

func unregisterHandler(w http.ResponseWriter, r *http.Request) {
	// incoming request format: /unregister/<service>
	logIt(r)
	service := strings.Split(r.URL.Path, "/")[2]
	delete(directory, service)
	directory.dump()
	messaging.Ok(w)
}

func generalHandler(w http.ResponseWriter, r *http.Request) {
	logIt(r)
	var err error
	message, err := messaging.GetMessage(r)
	if err == nil {
		service := message.Service
		if service == name && message.Command[0] == "get" && message.Command[2] == "version" {
			fmt.Fprint(w, version)
		} else {
			host, present := directory[service]
			if present {
				response, err := message.Send(host)
				if err != nil {
					messaging.FailInternal(w, err)
				} else {
					w.WriteHeader(response.Code)
					w.Header().Set("Content-Type", response.ContentType)
					fmt.Fprint(w, response.Payload)
				}
			} else if _, ok := services[service]; ok {
				messaging.FailInternal(w, fmt.Errorf("service '%v' is not active", service))
			} else {
				messaging.FailNotAuthorized(w)
			}
		}
	} else {
		messaging.FailInternal(w, err)
	}
}

func main() {
	var handlers = server.Handlers{
		{Pattern: "/register/", Handler: registerHandler},
		{Pattern: "/unregister/", Handler: unregisterHandler},
		{Pattern: "/", Handler: generalHandler},
	}
	me := server.NewServer(":"+port, &handlers)
	me.SetOnStart(func() { logger.GetLogger().Logf(logger.Info, "Listening on port %v", port) })
	me.SetOnFail(func(err error) { logger.GetLogger().Log(logger.Error, err) })
	me.SetOnInterrupt(func(signal os.Signal) { logger.GetLogger().Logf(logger.Warning, "Received %v", signal) })
	me.Start()
	logger.GetLogger().Log(logger.Info, "Exiting now")
}
