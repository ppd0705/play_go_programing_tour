package server

import (
	"context"
	pb "grpc-demo/proto"
	"io"
	"log"
)

type GreetSever struct {
}

var _ pb.GreeterServer = (*GreetSever)(nil)

func (s *GreetSever) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "hello.world"}, nil
}

func (s *GreetSever) SayList(r *pb.HelloRequest, stream pb.Greeter_SayListServer) error {
	for n := 0; n < 6; n++ {
		err := stream.Send(&pb.HelloReply{Message: "hello.list"})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *GreetSever) SayRecord(stream pb.Greeter_SayRecordServer) error {
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.HelloReply{Message: "say.record"})
		}
		if err != nil {
			return err
		}
		log.Printf("SayRecord resp: %v\n", resp)
	}
}

func (s *GreetSever) SayRoute(stream pb.Greeter_SayRouteServer) error {
	for {
		err := stream.Send(&pb.HelloReply{Message: "say.route"})
		if err != nil {
			return err
		}
		resp, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		log.Printf("resp: %v", resp)
	}
}
