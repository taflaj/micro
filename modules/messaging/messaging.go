// modules/messaging/messaging.go

package messaging

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/taflaj/micro/modules/logger"
)

const (
	version = "0.1.0 dev"
	agent   = "Microservices/" + version
	// Access is the name of the access service
	Access = "access"
)

// Message defines a message to be shared with other entities
type Message struct {
	To      []string
	CC      []string
	From    string
	Service string
	Method  string
	Request string
	Command []string
	IP      int
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

// Map is a map of strings
type Map map[string]string

// AccessLevelQuery allows a service to query on an entity's access level to a service
type AccessLevelQuery struct {
	Address int
	Service string
}

// AccessLevel contains an entity's access level to a service
type AccessLevel struct {
	Address  int
	Service  string
	Defined  bool
	Level    string
	CanRead  bool
	CanWrite bool
}

// Messenger is the default messenger
var Messenger *Host

func init() {
	Messenger = &Host{Address: "localhost", Port: "9998"}
}

func getResponse(response *http.Response, err error) (*Response, error) {
	if err == nil {
		defer response.Body.Close()
		// logger.GetLogger().Printf("%#v", response)
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
			// logger.GetLogger().Printf("%#v", result)
			return result, nil
		}
	}
	logger.GetLogger().Print(err)
	return nil, err
}

// Request performs a web service request
func Request(method string, url string, contentType string, body io.Reader, headers *Map) (*http.Response, error) {
	var (
		response *http.Response
		err      error
	)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err == nil {
		req.Header.Set("User-Agent", agent)
		req.Header.Set("Content-Type", contentType)
		if headers != nil {
			for k, v := range *headers {
				req.Header.Set(k, v)
			}
		}
		response, err = client.Do(req)
	}
	return response, err
}

// Send sends the message to a recipient
func (message *Message) Send(host *Host) (*Response, error) {
	// logger.GetLogger().Printf("Sending %#v to %#v", message, host)
	url := "http://" + host.Address + ":" + host.Port + "/"
	logger.GetLogger().Printf("  -> %v %#v", url, message)
	data, err := json.Marshal(message)
	// logger.GetLogger().Printf("DEBUG json data=%v", string(data))
	if err != nil {
		return nil, err
	}
	response, err := getResponse(Request(http.MethodPost, url, "application/json", bytes.NewBuffer(data), nil))
	logger.GetLogger().Printf("  <- %#v, %v", response, err)
	return response, err
}

// Get sends a request to a recipient
func Get(host *Host, request string) (*Response, error) {
	url := "http://" + host.Address + ":" + host.Port + "/" + request
	return getResponse(Request(http.MethodGet, url, "", nil, nil))
}

// GetMessage assembles a received message
func GetMessage(r *http.Request) (*Message, error) {
	var err error
	body, err := ioutil.ReadAll(r.Body)
	if err == nil {
		// logger.GetLogger().Printf("DEBUG json object=%v", string(body))
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

// QueryAccess queries an entity's access level to a service
func QueryAccess(address int, service string, from string) (*AccessLevel, error) {
	query := AccessLevelQuery{address, service}
	data, _ := json.Marshal(query)
	message := Message{
		To:      []string{Access},
		CC:      []string{},
		From:    from,
		Service: Access,
		Method:  http.MethodGet,
		Request: "",
		Command: []string{},
		IP:      -1,
		Data:    string(data),
	}
	response, err := message.Send(Messenger)
	logger.GetLogger().Printf("DEBUG response=%#v", response)
	// if err != nil {
	// 	return nil, err
	// }
	return nil, err
}

// Ok is good, not just okay
func Ok(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Ok")
}

// Fail returns an error message
func Fail(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	fmt.Fprint(w, err)
}

// FailBadRequest says it was a bad request
func FailBadRequest(w http.ResponseWriter, err error) {
	Fail(w, http.StatusBadRequest, err)
}

// FailBadRequestStandard says is was a bad request, but using a standard message
func FailBadRequestStandard(w http.ResponseWriter) {
	FailBadRequest(w, errors.New("Bad request format"))
}

// FailInternal oops, my bad
func FailInternal(w http.ResponseWriter, err error) {
	Fail(w, http.StatusInternalServerError, err)
}

// FailNotAuthorized says it was not authorized
func FailNotAuthorized(w http.ResponseWriter, err error) {
	Fail(w, http.StatusUnauthorized, err)
}

// FailNotAuthorizedStandard says it was not authorized, but using a standard message
func FailNotAuthorizedStandard(w http.ResponseWriter) {
	FailNotAuthorized(w, errors.New("Not authorized"))
}

// IPtoInt converts an IP address to integer
func IPtoInt(ip string) (int, error) {
	splits := strings.Split(ip, ".")
	if len(splits) != 4 {
		return -1, fmt.Errorf("'%v' is not a valid IP address", ip)
	}
	// logger.GetLogger().Printf("splits: %#v", splits)
	result := 0
	for _, split := range splits {
		// logger.GetLogger().Printf("split : %v", split)
		// result <<= 8
		b, err := strconv.ParseInt(split, 10, 9)
		// logger.GetLogger().Printf("err   : %v", err)
		if err != nil {
			return -1, err
		}
		// logger.GetLogger().Printf("b     : %v", b)
		result = result<<8 + int(b)
		// logger.GetLogger().Printf("result: %v", result)
	}
	return result, nil
}

// IPtoString converts an IP address to string
func IPtoString(ip int) string {
	return fmt.Sprintf("%v.%v.%v.%v", ip>>24, ip>>16&0xff, ip>>8&0xff, ip&0xff)
}

// Decompose converts data into a map
func (message *Message) Decompose() *Map {
	data := make(Map)
	// logger.GetLogger().Printf("DEBUG message.Data=%#v", message.Data)
	for _, value := range message.Data.([]interface{}) {
		splits := strings.Split(value.(string), "=")
		data[splits[0]] = splits[1]
	}
	return &data
}
