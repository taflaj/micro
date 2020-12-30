// server: main.go

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/taflaj/micro/messaging"
)

const (
	blank   = "<blank>"
	name    = "server"
	port    = "8888"
	version = "0.1.2"
)

var registered bool

func init() {
	log.SetFlags(log.Flags() | log.Lmicroseconds)
}

func logIt(r *http.Request) {
	agent := blank
	agents := r.Header["User-Agent"]
	if agents != nil && len(agents) > 0 {
		agent = agents[0]
	}
	log.Printf("%v %v from %v using %v", r.Method, r.URL.Path, r.RemoteAddr, agent)
}

func handler(w http.ResponseWriter, r *http.Request) {
	logIt(r)
	if !registered {
		if _, err := messaging.Get(messaging.Messenger, "register/"+name+"/"+port); err != nil {
			log.Print(err)
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
	log.Printf("  Contents: \"%v\" (%v); Accept \"%v\"", contentType, r.ContentLength, accept)
	ip := r.RemoteAddr
	p := strings.LastIndex(ip, ":")
	msg := &messaging.Message{From: name, Request: r.RequestURI, IP: ip[:p]}
	msg.Command = strings.Split(r.RequestURI[1:], "/")
	if len(msg.Command) > 1 {
		msg.Service = msg.Command[1]
	}
	msg.To = append(msg.To, msg.Service)
	// msg.CC = append(msg.CC, "logger")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
	}
	switch contentType {
	case "application/json":
		if err = json.Unmarshal(body, &msg.Data); err != nil {
			log.Print(err)
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
		// 		log.Printf("%v", msg.Data)
		// 	}
		// }
	}
	// switch r.Method {
	// case http.MethodDelete, http.MethodGet:
	// 	fmt.Fprintf(w, "Ok\n")
	// case http.MethodPost, http.MethodPut:
	// 	body, err := ioutil.ReadAll(r.Body)
	// 	if err != nil {
	// 		log.Printf("%#v", err)
	// 		w.WriteHeader(http.StatusBadRequest)
	// 	} else {
	// 		// if ct == "application/json" {}
	// 		w.WriteHeader(http.StatusOK)
	// 		fmt.Fprintf(w, "%v\n", string(body))
	// 	}
	// 	// r.ParseMultipartForm(1024)
	// 	// fmt.Fprintf(w, "Form: %#v\nPostForm: %#v\nMultipartForm: %#v\n", r.Form, r.PostForm, r.MultipartForm)
	// }
	log.Printf("  -> %#v", msg)
	response, err := msg.Send(messaging.Messenger)
	log.Printf("  <- %#v, %v\n", response, err)
	w.WriteHeader(response.Code)
	w.Header().Set("Content-Type", response.ContentType)
	fmt.Fprintf(w, "%v", response.Payload)
	// log.Printf("%#v", msg)
}

func main() {
	http.HandleFunc("/get/"+name+"/version", func(w http.ResponseWriter, r *http.Request) {
		logIt(r)
		// log.Printf("%v %v from %v using %v", r.Method, r.URL.Path, r.RemoteAddr, r.Header["User-Agent"][0])
		fmt.Fprint(w, version)
	})
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}