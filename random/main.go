// random: main.go

package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/taflaj/micro/modules/logger"
	"github.com/taflaj/micro/modules/messaging"
	pb "github.com/taflaj/micro/pubsub/pubsub"
	"github.com/taflaj/util/random"
	"google.golang.org/grpc"
)

const (
	host    = "localhost"
	name    = "random"
	port    = 9997
	version = "0.1.2"
)

type service struct {
	pb.UnimplementedPubSubServer
}

func init() {
	logger.NewLogger(name, logger.Info)
}

// Notify implements pubsub.Notify
func (s *service) Notify(ctx context.Context, in *pb.PublishMessage) (*pb.NotifyReply, error) {
	if len(in.Command) < 3 {
		return messaging.FailBadRequest(name)
	}
	length := 32
	if len(in.Command) > 3 {
		l, err := strconv.ParseInt(in.Command[3], 10, 0)
		if err == nil {
			length = int(l)
		}
	}
	var (
		result string
		err    error
	)
	t := in.Command[2]
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
		return messaging.FailBadRequest(name, fmt.Errorf("unknown argument \"%v\"", t))
	}
	if err == nil {
		return &pb.NotifyReply{From: name, Code: http.StatusOK, Data: result}, nil
	}
	return messaging.FailInternal(name, err)
}

var connection *grpc.ClientConn

func main() {
	var err error
	connection, err = grpc.Dial(messaging.Host, grpc.WithInsecure())
	if err == nil {
		defer connection.Close()
		client := pb.NewPubSubClient(connection)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		request := &pb.SubscribeRequest{Service: name, Host: host, Port: port}
		reply, err := client.Subscribe(ctx, request)
		if err == nil {
			id := reply.Id
			logger.GetLogger().Logf(logger.Info, "Subscribed with id %v", id)
			me := messaging.NewServer(fmt.Sprintf("%v:%v", host, port))
			pb.RegisterPubSubServer(me.GetServer(), &service{})
			me.SetOnStart(func(addr net.Addr) { logger.GetLogger().Logf(logger.Info, "Listening at %v", addr) })
			me.SetOnFail(func(err error) { logger.GetLogger().Log(logger.Critical, err) })
			me.SetOnInterrupt(func(signal os.Signal) { logger.GetLogger().Logf(logger.Warning, "Received signal: %v", signal) })
			me.ServeAndWait()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			reply, err := client.Unsubscribe(ctx, &pb.UnsubscribeRequest{Request: request, Id: id})
			if err == nil {
				logger.GetLogger().Logf(logger.Info, "Unsubscribed with id %v", reply.Id)
			} else {
				logger.GetLogger().Log(logger.Warning, err)
			}
		} else {
			logger.GetLogger().Log(logger.Critical, err)
		}
	} else {
		logger.GetLogger().Log(logger.Critical, err)
	}
	logger.GetLogger().Log(logger.Info, "Exiting now")
}
