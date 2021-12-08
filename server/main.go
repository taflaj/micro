// server: main.go

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/taflaj/micro/modules/logger"
	"github.com/taflaj/micro/modules/messaging"
	pb "github.com/taflaj/micro/pubsub/pubsub"
	"github.com/taflaj/util/ipinfo"
	"github.com/taflaj/util/server"
	"google.golang.org/grpc"
)

const (
	blank   = "<blank>"
	name    = "server"
	port    = 9999
	version = "0.1.5"
)

func init() {
	logger.NewLogger(name, logger.Info)
}

func getAddress(origin string) (string, uint32) {
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

var connection *grpc.ClientConn

func generalHandler(w http.ResponseWriter, r *http.Request) {
	logIt(r)
	contentType := getFirst(r.Header["Content-Type"])
	accept := getFirst(r.Header["Accept"])
	logger.GetLogger().Logf(logger.Info, "  Content-Type: \"%v\" (%v); Accept \"%v\"", contentType, r.ContentLength, accept)
	_, ip := getAddress(r.RemoteAddr)
	message := &pb.PublishMessage{From: name, Method: r.Method, Request: r.RequestURI, Ip: ip}
	message.Command = strings.Split(r.RequestURI[1:], "/")
	if len(message.Command) > 1 {
		service := message.Command[1]
		message.Service = service
		message.To = []string{service}
	}
	// logger.GetLogger().Logf(logger.Debug, "%#v", message)
	switch contentType {
	// case "application/json":
	// 	if err = json.Unmarshal(body, &message.Data); err != nil {
	// 		logger.GetLogger().Log(logger.Warning, err)
	// 	}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			logger.GetLogger().Log(logger.Warning, err)
		}
		b, err := json.Marshal(r.Form)
		if err == nil {
			message.Extra = string(b)
		} else {
			logger.GetLogger().Log(logger.Warning, err)
		}
	default:
		if strings.HasPrefix(contentType, "multipart/form-data") {
			if err := r.ParseMultipartForm(1048576); err != nil {
				logger.GetLogger().Log(logger.Warning, err)
			} else {
				b, err := json.Marshal(r.PostForm)
				if err == nil {
					message.Extra = string(b)
				} else {
					logger.GetLogger().Log(logger.Warning, err)
				}
			}
		} else {
			defer r.Body.Close()
			body, err := ioutil.ReadAll(r.Body)
			if err == nil {
				message.Extra = string(body)
			} else {
				logger.GetLogger().Log(logger.Warning, err)
			}
			// message.Extra = string(body) //url.Values{"": []string{string(body)}}
		}
	}
	client := pb.NewPubSubClient(connection)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	replies, err := client.Publish(ctx, message)
	if err == nil {
		logger.GetLogger().Logf(logger.Debug, "%#v", replies)
		first := replies.Replies[0]
		w.WriteHeader(int(first.Code))
		contentType = first.Type
		if len(contentType) == 0 {
			contentType = "text/plain"
		}
		w.Header().Set("Content-Type", contentType)
		fmt.Fprint(w, first.Data)
	} else {
		logger.GetLogger().Log(logger.Warning, err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
	}
	// fmt.Fprintf(w, "%v\n", message)
}

func main() {
	if len(os.Args) > 1 {
		ipinfo.SetToken(os.Args[1])
	}
	var err error
	connection, err = grpc.Dial(messaging.Host, grpc.WithInsecure())
	if err == nil {
		defer connection.Close()
		var handlers = server.Handlers{
			{Pattern: "/get/" + name + "/version", Handler: versionHandler},
			{Pattern: "/", Handler: generalHandler},
		}
		me := server.NewServer(fmt.Sprintf(":%v", port), &handlers)
		me.SetOnStart(func() { logger.GetLogger().Logf(logger.Info, "Listening on port %v", port) })
		me.SetOnFail(func(err error) { logger.GetLogger().Log(logger.Error, err) })
		me.SetOnInterrupt(func(signal os.Signal) { logger.GetLogger().Logf(logger.Warning, "Received signal: %v", signal) })
		me.Start()
	} else {
		logger.GetLogger().Log(logger.Critical, err)
	}
	logger.GetLogger().Log(logger.Info, "Exiting now")
}
