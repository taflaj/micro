// server: main.go

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/taflaj/micro/modules/logger"
	"github.com/taflaj/micro/modules/messaging"
)

const (
	blank   = "<blank>"
	name    = "server"
	port    = "9999"
	version = "0.1.4 dev"
)

var (
	registered bool
	// token      string
	headers messaging.Map
	// log     logger.Logger
)

func init() {
	headers = messaging.Map{"Accept": "application/json"}
	logger.NewLogger(name)
}

func getIP(ip string) (string, int) {
	p := strings.LastIndex(ip, ":")
	asString := ip[:p]
	asInt, _ := messaging.IPtoInt(ip[:p])
	return asString, asInt
}

func logIt(r *http.Request) {
	agent := blank
	agents := r.Header["User-Agent"]
	if agents != nil && len(agents) > 0 {
		agent = agents[0]
	}
	// logger.GetLogger().Printf("%#v", r)
	logger.GetLogger().Printf("%v %v from %v using %v", r.Method, r.URL.Path, r.RemoteAddr, agent)
	// logger.GetLogger().Spy(agent)
	go func(r *http.Request) {
		// var headers *messaging.Header
		// headers := messaging.Header{}
		ip, _ := getIP(r.RemoteAddr)
		url := "http://ipinfo.io/" + ip
		// if token != "" {
		// 	// url += "?token=" + token
		// 	// headers = &messaging.Header{"Authorization": "Bearer " + token}
		// 	headers["Authorization"] = "Bearer " + token
		// }
		var err error
		response, err := messaging.Request(http.MethodGet, url, "", nil, &headers)
		// logger.GetLogger().Printf("%#v", response)
		if err == nil {
			defer response.Body.Close()
			data, err := ioutil.ReadAll(response.Body)
			if err == nil {
				// logger.GetLogger().Printf("%v", string(data))
				var ipinfo struct {
					IP       string
					City     string
					Region   string
					Country  string
					Postal   string
					HostName string
					Org      string
					Bogon    bool
				}
				if err = json.Unmarshal(data, &ipinfo); err == nil {
					// logger.GetLogger().Printf("%&v", ipinfo)
					report := fmt.Sprintf("  From: %v", ipinfo.IP)
					if !ipinfo.Bogon {
						report = fmt.Sprintf("%v %v/%v/%v/%v|%v|%v", report, ipinfo.City, ipinfo.Region, ipinfo.Postal, ipinfo.Country, ipinfo.HostName, ipinfo.Org)
					}
					logger.GetLogger().Print(report)
					// if ipinfo.Bogon {
					// 	logger.GetLogger().Printf("  From: %v", ipinfo.IP)
					// } else {
					// 	logger.GetLogger().Printf("  From: %v %v/%v/%v/%v|%v|%v")
					// }
					// logger.GetLogger().Printf("  %#v", ipinfo)
					return
				}
			}
		}
		logger.GetLogger().Print(err)
	}(r)
}

func handler(w http.ResponseWriter, r *http.Request) {
	logIt(r)
	if !registered {
		if _, err := messaging.Get(messaging.Messenger, "register/"+name+"/"+port+"/localhost"); err != nil {
			logger.GetLogger().Print(err)
		} else {
			registered = true
		}
	}
	contentType := ""
	if r.ContentLength > 0 {
		contentType = r.Header["Content-Type"][0]
	}
	accept := blank
	accepts := r.Header["Accept"]
	if accepts != nil && len(accepts) > 0 {
		accept = accepts[0]
	}
	logger.GetLogger().Printf("  Contents: \"%v\" (%v); Accept \"%v\"", contentType, r.ContentLength, accept)
	_, ip := getIP((r.RemoteAddr))
	msg := &messaging.Message{From: name, Method: r.Method, Request: r.RequestURI, IP: ip}
	msg.Command = strings.Split(r.RequestURI[1:], "/")
	if len(msg.Command) > 1 {
		msg.Service = msg.Command[1]
	}
	msg.To = append(msg.To, msg.Service)
	// msg.CC = append(msg.CC, "logger")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.GetLogger().Print(err)
	}
	switch contentType {
	case "application/json":
		if err = json.Unmarshal(body, &msg.Data); err != nil {
			logger.GetLogger().Print(err)
		}
	case "application/x-www-form-urlencoded":
		msg.Data = strings.Split(string(body), "&")
	default:
		msg.Data = body
		// splits := strings.Split(contentType, ";")
		// if len(splits) > 0 {
		// 	if splits[0] == "multipart/form-data" {
		// 		boundary := splits[1][10:]
		// 		msg.Data = strings.Split(string(body), boundary)
		// 		logger.GetLogger().Printf("%v", msg.Data)
		// 	}
		// }
	}
	// switch r.Method {
	// case http.MethodDelete, http.MethodGet:
	// 	fmt.Fprintf(w, "Ok\n")
	// case http.MethodPost, http.MethodPut:
	// 	body, err := ioutil.ReadAll(r.Body)
	// 	if err != nil {
	// 		logger.GetLogger().Printf("%#v", err)
	// 		w.WriteHeader(http.StatusBadRequest)
	// 	} else {
	// 		// if ct == "application/json" {}
	// 		w.WriteHeader(http.StatusOK)
	// 		fmt.Fprintf(w, "%v\n", string(body))
	// 	}
	// 	// r.ParseMultipartForm(1024)
	// 	// fmt.Fprintf(w, "Form: %#v\nPostForm: %#v\nMultipartForm: %#v\n", r.Form, r.PostForm, r.MultipartForm)
	// }
	// if msg.Service == "" {
	// 	messaging.FailBadRequest(w)
	// 	logger.GetLogger().Print("  Invalid command")
	// } else {
	// logger.GetLogger().Printf("  -> %#v", msg)
	response, err := msg.Send(messaging.Messenger)
	// logger.GetLogger().Printf("  <- %#v, %v", response, err)
	w.WriteHeader(response.Code)
	w.Header().Set("Content-Type", response.ContentType)
	fmt.Fprintf(w, "%v", response.Payload)
	// }
	// logger.GetLogger().Printf("%#v", msg)
}

func main() {
	if len(os.Args) > 1 {
		// token = os.Args[1]
		headers["Authorization"] = "Bearer " + os.Args[1]
	}
	http.HandleFunc("/get/"+name+"/version", func(w http.ResponseWriter, r *http.Request) {
		logIt(r)
		// logger.GetLogger().Printf("%v %v from %v using %v", r.Method, r.URL.Path, r.RemoteAddr, r.Header["User-Agent"][0])
		fmt.Fprint(w, version)
	})
	http.HandleFunc("/", handler)
	logger.GetLogger().Printf("Listening on port %v", port)
	logger.GetLogger().Fatal(http.ListenAndServe(":"+port, nil))
}
