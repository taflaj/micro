// modules/messaging/messaging.go

package messaging

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	pb "github.com/taflaj/micro/pubsub/pubsub"
	"google.golang.org/grpc"
)

const (
	Host = "localhost:9998"
)

type Server interface {
	GetServer() *grpc.Server
	GracefulStop()
	SetOnExit(func())
	SetOnFail(func(error))
	SetOnInterrupt(func(os.Signal))
	SetOnStart(func(net.Addr))
	ServeAndWait()
	Stop()
}

type serverSetup struct {
	server      *grpc.Server
	address     string
	onStart     func(net.Addr)
	onExit      func()
	onFail      func(error)
	onInterrupt func(os.Signal)
}

func NewServer(address string) Server {
	server := serverSetup{}
	server.server = grpc.NewServer()
	server.address = address
	return &server
}

func (ss *serverSetup) ServeAndWait() {
	go func() {
		listener, err := net.Listen("tcp", ss.address)
		if err == nil {
			if ss.onStart != nil {
				ss.onStart(listener.Addr())
			}
			if err := ss.server.Serve(listener); err != nil {
				if ss.onFail != nil {
					ss.onFail(err)
				}
			}
		} else if ss.onFail != nil {
			ss.onFail(err)
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGUSR1)
	sig := <-stop
	if ss.onInterrupt != nil {
		ss.onInterrupt(sig)
	}
	ss.server.GracefulStop()
	if ss.onExit != nil {
		ss.onExit()
	}
}

func (ss *serverSetup) exit(graceful bool) {
	if graceful {
		ss.server.GracefulStop()
	} else {
		ss.server.Stop()
	}
	if ss.onExit != nil {
		ss.onExit()
	}
}

func (ss *serverSetup) GracefulStop() {
	ss.exit(true)
}

func (ss *serverSetup) Stop() {
	ss.exit(false)
}

func (ss *serverSetup) GetServer() *grpc.Server {
	return ss.server
}

func (ss *serverSetup) SetOnStart(f func(net.Addr)) {
	ss.onStart = f
}

func (ss *serverSetup) SetOnFail(f func(error)) {
	ss.onFail = f
}

func (ss *serverSetup) SetOnInterrupt(f func(os.Signal)) {
	ss.onInterrupt = f
}

func (ss *serverSetup) SetOnExit(f func()) {
	ss.onExit = f
}

// IPtoInt converts an IP address to integer
func IPtoInt(ip string) (uint32, error) {
	splits := strings.Split(ip, ".")
	if len(splits) != 4 {
		return 0, fmt.Errorf("'%v' is not a valid IP address", ip)
	}
	var result uint32 = 0
	for _, split := range splits {
		b, err := strconv.ParseInt(split, 10, 9)
		if err != nil {
			return 0, err
		}
		result = result<<8 + uint32(b)
	}
	return result, nil
}

// IPtoString converts an IP address to string
func IPtoString(ip int) string {
	return fmt.Sprintf("%v.%v.%v.%v", ip>>24, ip>>16&0xff, ip>>8&0xff, ip&0xff)
}

// Fail returns an error message
func Fail(name string, code uint32, err error) (*pb.NotifyReply, error) {
	reply := &pb.NotifyReply{From: name, Code: code, Data: err.Error()}
	return reply, err
}

// FailBadRequest says it was a bad request
func FailBadRequest(name string, err ...error) (*pb.NotifyReply, error) {
	e := errors.New("bad request format")
	if len(err) > 0 {
		e = err[0]
	}
	return Fail(name, http.StatusBadRequest, e)
}

// FailInternal oops, my bad
func FailInternal(name string, err error) (*pb.NotifyReply, error) {
	return Fail(name, http.StatusInternalServerError, err)
}

// FailUnavailable says the service is not available
func FailUnavailable(name string, err error) (*pb.NotifyReply, error) {
	return Fail(name, http.StatusServiceUnavailable, err)
}
