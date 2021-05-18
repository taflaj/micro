// modules/messaging/messaging_test.go

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

	"github.com/taflaj/services/modules/logger"
)

const (
	version = "0.1.0"
	agent   = "Microservices/" + version
	// Access is the name of the access service
	Access = "access"
)

// Map is a map of strings
type Map map[string]string

// Host specifies a recipient
type Host struct {
	Address string
	Port    string
}

// Message defines a message to be shared with other entities
type Message struct {
	From    string
	Service string
	Method  string
	Request string
	Command []string
	IP      int
	Data    []string
}

// Response contains an http response
type Response struct {
	Code        int
	Length      int64
	ContentType string
	Payload     string
}

// AccessLevel contains an entity's access level to a service
type AccessLevel struct {
	IP       int
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
		var data []byte
		defer response.Body.Close()
		logger.GetLogger().Logf(logger.Debug, "%#v", response)
		data, err = ioutil.ReadAll(response.Body)
		if err == nil {
			result := Response{
				Code:    response.StatusCode,
				Length:  response.ContentLength,
				Payload: string(data),
			}
			contentType, ok := response.Header["Content-Type"]
			if ok && len(contentType) > 0 {
				result.ContentType = contentType[0]
			}
			logger.GetLogger().Logf(logger.Debug, "%#v", result)
			return &result, nil
		}
	}
	return nil, err
}

// Request performs a web service request
func Request(method string, url string, contentType string, body io.Reader, headers *Map) (*http.Response, error) {
	success := make(chan *http.Response)
	failure := make(chan error)
	go func() {
		client := http.Client{}
		request, err := http.NewRequest(method, url, body)
		if err != nil {
			failure <- err
		}
		request.Header.Set("User-Agent", agent)
		request.Header.Set("Content-Type", contentType)
		if headers != nil {
			for k, v := range *headers {
				request.Header.Set(k, v)
			}
		}
		response, err := client.Do(request)
		if err != nil {
			failure <- err
		}
		success <- response
		logger.GetLogger().Logf(logger.Debug, "%#v", response)
	}()
	select {
	case response := <-success:
		return response, nil
	case err := <-failure:
		return nil, err
	}
}

// Send sends a message to a recipient
func (message *Message) Send(host *Host) (*Response, error) {
	url := "http://" + host.Address + ":" + host.Port + "/"
	logger.GetLogger().Logf(logger.Info, "  -> %v %#v", url, message)
	data, err := json.Marshal(message)
	if err != nil {
		logger.GetLogger().Log(logger.Error, err)
		return nil, err
	}
	response, err := getResponse(Request(http.MethodPost, url, "application/json", bytes.NewBuffer(data), nil))
	logger.GetLogger().Logf(logger.Info, "  <- %#v, %v", response, err)
	return response, err
}

// GetData converts the data in the Message into a Map
func (message *Message) GetData() *Map {
	data := make(Map)
	for _, input := range message.Data {
		splits := strings.Split(input, "=")
		value := ""
		if len(splits) == 2 {
			value = splits[1]
		}
		data[splits[0]] = value
	}
	return &data
}

// Get sends a request to a recipient
func Get(host *Host, request string) (*Response, error) {
	url := "http://" + host.Address + ":" + host.Port
	if request[0] != '/' {
		url += "/"
	}
	url += request
	return getResponse(Request(http.MethodGet, url, "", nil, nil))
}

// GetMessage assembles a received message
func GetMessage(r *http.Request) (*Message, error) {
	var err error
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err == nil {
		err = errors.New("invalid message format")
		if r.ContentLength > 0 && r.Header["Content-Type"][0] == "application/json" {
			message := Message{}
			if err = json.Unmarshal(body, &message); err == nil {
				return &message, nil
			}
		}
	}
	return nil, err
}

// Ok is good, not just okay
func Ok(w http.ResponseWriter) {
	fmt.Fprint(w, "Ok")
}

// Fail returns an error message
func Fail(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	fmt.Fprint(w, err)
}

// FailBadRequest says it was a bad request
func FailBadRequest(w http.ResponseWriter, err ...error) {
	e := errors.New("bad request format")
	if len(err) > 0 {
		e = err[0]
	}
	Fail(w, http.StatusBadRequest, e)
}

// FailInternal oops, my bad
func FailInternal(w http.ResponseWriter, err error) {
	Fail(w, http.StatusInternalServerError, err)
}

// FailNotAuthorized says it was not authorized
func FailNotAuthorized(w http.ResponseWriter) {
	Fail(w, http.StatusUnauthorized, errors.New("not authorized"))
}

// IPtoInt converts an IP address to integer
func IPtoInt(ip string) (int, error) {
	splits := strings.Split(ip, ".")
	if len(splits) != 4 {
		return -1, fmt.Errorf("'%v' is not a valid IP address", ip)
	}
	result := 0
	for _, split := range splits {
		b, err := strconv.ParseInt(split, 10, 9)
		if err != nil {
			return -1, err
		}
		result = result<<8 + int(b)
	}
	return result, nil
}

// IPtoString converts an IP address to string
func IPtoString(ip int) string {
	return fmt.Sprintf("%v.%v.%v.%v", ip>>24, ip>>16&0xff, ip>>8&0xff, ip&0xff)
}
