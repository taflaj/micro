// pubsub: main.go

package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/taflaj/micro/modules/logger"
	"github.com/taflaj/micro/modules/messaging"
	pb "github.com/taflaj/micro/pubsub/pubsub"
	"google.golang.org/grpc"
)

const (
	name    = "pubsub"
	version = "0.1.4"
)

type service struct {
	pb.UnimplementedPubSubServer
}

type subscriber struct {
	host       string
	port       uint16
	active     bool
	connection *grpc.ClientConn
}

type subscribers map[uint64]*subscriber

var channels map[string]subscribers

func init() {
	logger.NewLogger(name, logger.Info)
	channels = make(map[string]subscribers)
}

// Subscribe implements pubsub.Subscribe
func (s *service) Subscribe(ctx context.Context, in *pb.SubscribeRequest) (*pb.SubscribeReply, error) {
	logger.GetLogger().Logf(logger.Debug, "Received %#v", in)
	logger.GetLogger().Logf(logger.Info, "Adding %v:%v to channel %v", in.Host, in.Port, in.Service)
	connection, err := grpc.Dial(fmt.Sprintf("%v:%v", in.Host, in.Port), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	sub := &subscriber{host: in.Host, port: uint16(in.Port), active: true, connection: connection}
	id := uint64(reflect.ValueOf(sub).Pointer())
	_, ok := channels[in.Service]
	if !ok {
		channels[in.Service] = make(subscribers)
	}
	channels[in.Service][id] = sub
	logger.GetLogger().Logf(logger.Debug, "%v", channels)
	for k1, v1 := range channels {
		logger.GetLogger().Logf(logger.Debug, "  %v: %v", k1, v1)
		for k2, v2 := range v1 {
			logger.GetLogger().Logf(logger.Debug, "    %v: %v", k2, v2)
		}
	}
	return &pb.SubscribeReply{Id: id}, nil
}

// Unsubscribe implements pubsub.Unsubscribe
func (s *service) Unsubscribe(ctx context.Context, in *pb.UnsubscribeRequest) (*pb.SubscribeReply, error) {
	logger.GetLogger().Logf(logger.Info, "Removing %v:%v as %v from channel %v", in.Request.Host, in.Request.Port, in.Id, in.Request.Service)
	channel, ok := channels[in.Request.Service]
	if !ok {
		return nil, fmt.Errorf("unknown channel %v", in.Request.Service)
	}
	sub, ok := channel[in.Id]
	if !ok {
		return nil, fmt.Errorf("invalid id %v for channel %v", in.Id, in.Request.Service)
	}
	if sub.host == in.Request.Host && sub.port == uint16(in.Request.Port) {
		if sub.active {
			sub.active = false
			sub.connection.Close()
			logger.GetLogger().Logf(logger.Debug, "%v", channels)
			for k1, v1 := range channels {
				logger.GetLogger().Logf(logger.Debug, "  %v: %v", k1, v1)
				for k2, v2 := range v1 {
					logger.GetLogger().Logf(logger.Debug, "    %v: %v", k2, v2)
				}
			}
			return &pb.SubscribeReply{Id: in.Id}, nil
		}
		return nil, fmt.Errorf("id %v was already unsubscribed from channel %v", in.Id, in.Request.Service)
	}
	return nil, fmt.Errorf("id %v does not match subscription", in.Id)
}

func notify(who string, message *pb.PublishMessage) []chan *pb.NotifyReply {
	logger.GetLogger().Logf(logger.Debug, "Notifying subscribers to channel %v", who)
	var replies []chan *pb.NotifyReply
	channel, ok := channels[who]
	if ok {
		for _, sub := range channel {
			reply := make(chan *pb.NotifyReply)
			replies = append(replies, reply)
			go func(sub *subscriber) {
				client := pb.NewPubSubClient(sub.connection)
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				r, err := client.Notify(ctx, message)
				if err != nil {
					r, _ = messaging.Fail(who, http.StatusInternalServerError, err)
				}
				reply <- r
			}(sub)
		}
	} else {
		logger.GetLogger().Logf(logger.Warning, "Service %v is unavailable", who)
		reply := make(chan *pb.NotifyReply)
		replies = append(replies, reply)
		go func() {
			r, _ := messaging.FailUnavailable(name, fmt.Errorf("service %v is unavailable", who))
			reply <- r
		}()
	}
	return replies
}

// Publish implements pubsub.Publish
func (s *service) Publish(ctx context.Context, in *pb.PublishMessage) (*pb.PublishReply, error) {
	logger.GetLogger().Logf(logger.Info, "%#v", in)
	if in.Service == name {
		var myReply *pb.NotifyReply
		if in.Command[0] == "get" && in.Command[2] == "version" {
			myReply = &pb.NotifyReply{From: name, Code: http.StatusOK, Data: version}
		} else {
			myReply, _ = messaging.FailBadRequest(name)
		}
		return &pb.PublishReply{Replies: []*pb.NotifyReply{myReply}}, nil
	}
	var replies []*pb.NotifyReply
	if len(in.To) == 0 {
		myReply, _ := messaging.FailBadRequest(name)
		replies = append(replies, myReply)
	} else {
		var responses []chan *pb.NotifyReply
		for _, sub := range in.To {
			responses = append(responses, notify(sub, in)...)
		}
		for _, reply := range responses {
			replies = append(replies, <-reply)
		}
	}
	return &pb.PublishReply{Replies: replies}, nil
}

func main() {
	me := messaging.NewServer(messaging.Host)
	pb.RegisterPubSubServer(me.GetServer(), &service{})
	me.SetOnStart(func(addr net.Addr) { logger.GetLogger().Logf(logger.Info, "Listening at %v", addr) })
	me.SetOnFail(func(err error) { logger.GetLogger().Log(logger.Critical, err) })
	me.SetOnInterrupt(func(signal os.Signal) { logger.GetLogger().Logf(logger.Warning, "Received signal: %v", signal) })
	me.ServeAndWait()
	logger.GetLogger().Log(logger.Info, "Exiting now")
}
