// messaging/messaging.go

package messaging

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Message defines a message to be shared with other entities
type Message struct {
	To      []string
	CC      []string
	From    string
	Service string
	Request string
	Command []string
	Data    interface{}
}

// Response contains an http response
type Response struct {
	Code        int
	Length      int64
	ContentType string
	Payload     interface{}
}

// Host specifies a recipient
type Host struct {
	Address string
	Port    string
}

// Messenger is the default messenger
var Messenger *Host

func init() {
	Messenger = &Host{Address: "localhost", Port: "8001"}
}

func getResponse(response *http.Response, err error) (*Response, error) {
	if err == nil {
		defer response.Body.Close()
		// log.Printf("%#v", response)
		data, err := ioutil.ReadAll(response.Body)
		if err == nil {
			result := &Response{
				Code:        response.StatusCode,
				Length:      response.ContentLength,
				ContentType: response.Header["Content-Type"][0],
				Payload:     data,
			}
			ct := strings.Split(result.ContentType, ";")[0]
			if ct == "text/plain" {
				result.Payload = string(data)
			}
			// log.Printf("%#v", result)
			return result, nil
		}
	}
	log.Print(err)
	return nil, err
}

// Send sends the message to a recipient
func (message *Message) Send(host *Host) (*Response, error) {
	// log.Printf("Sending %#v to %#v", message, host)
	url := "http://" + host.Address + ":" + host.Port + "/"
	data, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	return getResponse(http.Post(url, "application/json", bytes.NewBuffer(data)))
}

// Get sends a request to a recipient
func Get(host *Host, request string) (*Response, error) {
	url := "http://" + host.Address + ":" + host.Port + "/" + request
	return getResponse(http.Get(url))
}

// GetMessage assembles a received message
func GetMessage(r *http.Request) (*Message, error) {
	var err error
	body, err := ioutil.ReadAll(r.Body)
	if err == nil {
		// log.Print(string(body))
		err = errors.New("Invalid message format")
		if r.ContentLength > 0 && r.Header["Content-Type"][0] == "application/json" {
			msg := Message{}
			if err = json.Unmarshal(body, &msg); err == nil {
				return &msg, nil
			}
		}
	}
	return nil, err
}
