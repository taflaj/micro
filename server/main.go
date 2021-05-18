// server: main.go

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/taflaj/services/modules/logger"
	"github.com/taflaj/services/modules/messaging"
	"github.com/taflaj/util/ipinfo"
	"github.com/taflaj/util/server"
)

const (
	blank   = "<blank>"
	name    = "server"
	port    = "9999"
	version = "0.1.4"
)

func init() {
	logger.NewLogger(name, logger.Info)
	// logger.NewLogger(name, logger.Debug)
}

func getAddress(origin string) (string, int) {
	p := strings.LastIndex(origin, ":")
	asString := origin[:p]
	asInt, _ := messaging.IPtoInt(asString)
	return asString, asInt
}

func logIt(r *http.Request) {
	agent := blank
	agents := r.Header["User-Agent"]
	if len(agents) > 0 {
		agent = agents[0]
	}
	logger.GetLogger().Logdf(logger.Debug, 1, "%#v", r)
	logger.GetLogger().Logdf(logger.Info, 1, "%v %v from %v using %v", r.Method, r.URL.Path, r.RemoteAddr, agent)
	go func() {
		origin, _ := getAddress(r.RemoteAddr)
		info, err := ipinfo.GetInfo(origin)
		var output string
		if err == nil {
			if info.Bogon {
				output = info.IP
			} else {
				output = fmt.Sprintf("%s %s/%s/%s/%s|%s|%s", info.IP, info.City, info.Region, info.Postal, info.Country, info.HostName, info.Org)
			}
		} else {
			output = origin
		}
		logger.GetLogger().Log(logger.Info, "  From: "+output)
	}()
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	logIt(r)
	fmt.Fprint(w, version)
}

func getFirst(v []string) string {
	var result string
	if len(v) > 0 {
		result = v[0]
	}
	return result
}

func generalHandler(w http.ResponseWriter, r *http.Request) {
	logIt(r)
	contentType := getFirst(r.Header["Content-Type"])
	accept := getFirst(r.Header["Accept"])
	logger.GetLogger().Logf(logger.Info, "  Contents: \"%v\" (%v); Accept \"%v\"", contentType, r.ContentLength, accept)
	_, ip := getAddress(r.RemoteAddr)
	message := messaging.Message{From: name, Method: r.Method, Request: r.RequestURI, IP: ip}
	message.Command = strings.Split(r.RequestURI[1:], "/")
	if len(message.Command) > 1 {
		message.Service = message.Command[1]
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.GetLogger().Log(logger.Warning, err)
	}
	switch contentType {
	case "application/json":
		if err = json.Unmarshal(body, &message.Data); err != nil {
			logger.GetLogger().Log(logger.Warning, err)
		}
	case "application/x-www-form-urlencoded":
		message.Data = strings.Split(string(body), "&")
	default:
		message.Data = append(message.Data, string(body))
	}
	logger.GetLogger().Logf(logger.Debug, "%#v", message)
	response, err := message.Send(messaging.Messenger)
	if err != nil {
		logger.GetLogger().Log(logger.Warning, err)
	}
	w.WriteHeader(response.Code)
	w.Header().Set("Content-Type", response.ContentType)
	fmt.Fprintf(w, "%v", response.Payload)
}

func main() {
	if len(os.Args) > 1 {
		ipinfo.SetToken(os.Args[1])
	}
	response, err := messaging.Get(messaging.Messenger, "register/"+name+"/"+port+"/localhost")
	if err != nil {
		logger.GetLogger().Log(logger.Warning, err)
	} else if response.Code != http.StatusOK {
		logger.GetLogger().Logf(logger.Warning, "Server registration code: %v", response.Code)
	}
	var handlers = server.Handlers{
		{Pattern: "/get/" + name + "/version", Handler: versionHandler},
		{Pattern: "/", Handler: generalHandler},
	}
	me := server.NewServer(":"+port, &handlers)
	me.SetOnStart(func() { logger.GetLogger().Logf(logger.Info, "Listening on port %v", port) })
	me.SetOnFail(func(err error) { logger.GetLogger().Log(logger.Error, err) })
	me.SetOnInterrupt(func(signal os.Signal) { logger.GetLogger().Logf(logger.Warning, "Received %v", signal) })
	me.Start()
	messaging.Get(messaging.Messenger, "unregister/"+name)
	logger.GetLogger().Log(logger.Info, "Exiting now")
}
