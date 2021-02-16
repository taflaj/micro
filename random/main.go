// random: main.go

package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/taflaj/micro/modules/logger"
	"github.com/taflaj/micro/modules/messaging"
	"github.com/taflaj/util/random"
)

const (
	name    = "random"
	port    = "9997"
	version = "0.1.1 dev"
)

var log logger.Logger

func init() {
	logger.NewLogger(name)
}

func logIt(r *http.Request) {
	who := r.Header["User-Agent"][0]
	logger.GetLogger().Printf("%v %v from %v using %v", r.Method, r.URL.Path, r.RemoteAddr, who)
	// logger.GetLogger().Spy(who)
}

func handler(w http.ResponseWriter, r *http.Request) {
	logIt(r)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	msg, err := messaging.GetMessage(r)
	if err != nil {
		messaging.Fail(w, http.StatusInternalServerError, err)
	} else if len(msg.Command) < 3 {
		messaging.FailBadRequestStandard(w)
	} else {
		// logger.GetLogger().Printf("%#v", msg)
		length := 32
		if len(msg.Command) > 3 {
			l, err := strconv.ParseInt(msg.Command[3], 10, 64)
			if err == nil {
				length = int(l)
			}
		}
		var result string
		t := msg.Command[2]
		switch t {
		case "alpha":
			result, err = random.Alpha(length)
		case "alphanum":
			result, err = random.AlphaNum(length)
		case "any":
			result, err = random.Any(length)
		case "hex":
			result, err = random.Hex(length)
		case "number":
			result, err = random.Number(length)
		case "special":
			result, err = random.Special(length)
		case "version":
			result = version
		default:
			err = fmt.Errorf("Unknown argument \"%v\"", t)
		}
		if err == nil {
			fmt.Fprint(w, result)
		} else {
			messaging.Fail(w, http.StatusInternalServerError, err)
		}
	}
}

func main() {
	http.HandleFunc("/", handler)
	go func() {
		if _, err := messaging.Get(messaging.Messenger, "register/"+name+"/"+port+"/localhost"); err != nil {
			logger.GetLogger().Fatal(err)
		}
	}()
	logger.GetLogger().Printf("Listening on port %v", port)
	logger.GetLogger().Fatal(http.ListenAndServe(":"+port, nil))
}
