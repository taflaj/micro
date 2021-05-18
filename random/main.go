// ramdom: main.go

package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/taflaj/services/modules/logger"
	"github.com/taflaj/services/modules/messaging"
	"github.com/taflaj/util/random"
	"github.com/taflaj/util/server"
)

const (
	name    = "random"
	port    = "9997"
	version = "0.1.1"
)

func init() {
	logger.NewLogger(name, logger.Info)
}

func logIt(r *http.Request) {
	who := r.Header["User-Agent"][0]
	logger.GetLogger().Logdf(logger.Debug, 1, "%#v", r)
	logger.GetLogger().Logdf(logger.Info, 1, "%v %v from %v using %v", r.Method, r.URL.Path, r.RemoteAddr, who)
}

func handler(w http.ResponseWriter, r *http.Request) {
	logIt(r)
	message, err := messaging.GetMessage(r)
	if err != nil {
		messaging.FailInternal(w, err)
	} else if len(message.Command) < 3 {
		messaging.FailBadRequest(w)
	} else {
		length := 32
		if len(message.Command) > 3 {
			l, err := strconv.ParseInt(message.Command[3], 10, 0)
			if err == nil {
				length = int(l)
			}
		}
		var result string
		t := message.Command[2]
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
			messaging.FailBadRequest(w, fmt.Errorf("unknown argument \"%v\"", t))
			return
		}
		if err == nil {
			fmt.Fprint(w, result)
		} else {
			messaging.FailInternal(w, err)
		}
	}
}

func main() {
	messaging.Get(messaging.Messenger, "register/"+name+"/"+port+"/localhost")
	var handler = server.Handlers{{Pattern: "/", Handler: handler}}
	me := server.NewServer(":"+port, &handler)
	me.SetOnStart(func() { logger.GetLogger().Logf(logger.Info, "Listening on port %v", port) })
	me.SetOnFail(func(err error) { logger.GetLogger().Log(logger.Error, err) })
	me.SetOnInterrupt(func(signal os.Signal) { logger.GetLogger().Logf(logger.Warning, "Received %v", signal) })
	me.Start()
	logger.GetLogger().Log(logger.Info, "Exiting now")
}
