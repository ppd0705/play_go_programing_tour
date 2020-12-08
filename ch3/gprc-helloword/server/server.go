package server

import (
	"context"
	pb "grpc-demo/proto"
)

type GreetSever struct {
}

func (s *GreetSever) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "hello.world"}, nil
}

