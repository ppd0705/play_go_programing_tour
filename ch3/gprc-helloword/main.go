package main

import (
	"context"
	"google.golang.org/grpc"
	pb "grpc-demo/proto"
	"grpc-demo/server"
	"io"
	"log"
	"net"
)

func startServer(portCH chan string) {
	svc := grpc.NewServer()
	pb.RegisterGreeterServer(svc, &server.GreetSever{})
	lis, _ := net.Listen("tcp", ":0")
	portCH <- lis.Addr().String()
	log.Printf("start server %s\n", lis.Addr().String())
	_ = svc.Serve(lis)
}

func SayHello(client pb.GreeterClient) error {
	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "ppd"})
	if err != nil {
		return err
	}
	log.Printf("client Say hello resp: %s", resp.Message)
	return nil
}

func SayList(client pb.GreeterClient) error {
	stream, err := client.SayList(context.Background(), &pb.HelloRequest{Name: "ppd"})
	if err != nil {
		return err
	}
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		log.Printf("SayList strema recv: %v\n", resp)
	}
}

func SayRecord(client pb.GreeterClient) error {
	stream, err := client.SayRecord(context.Background())
	if err != nil {
		return err
	}
	for i := 0; i < 6; i++ {
		err := stream.Send(&pb.HelloRequest{Name: "ppd"})
		if err != nil {
			return err
		}
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	log.Printf("resp err: %v", resp)
	return nil
}

func startClient(port string) {
	conn, _ := grpc.Dial(port, grpc.WithInsecure())
	defer conn.Close()

	client := pb.NewGreeterClient(conn)
	//err := SayHello(client)
	//err := SayList(client)
	//err := SayRecord(client)
	err := SayRoute(client)
	if err != nil {
		log.Fatalf("Say hello err: %v", err)
	}
}

func SayRoute(client pb.GreeterClient) error {
	stream, err := client.SayRoute(context.Background())
	if err != nil {
		return err
	}
	for n := 0;n < 6;n++ {
		err = stream.Send(&pb.HelloRequest{Name: "ppd"})
		if err != nil {
			return err
		}
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Printf("Sayroute resp: %v", resp)
	}
	err = stream.CloseSend()
	if err != nil {
		return err
	}
	return nil
}
func main() {
	portCH := make(chan string)
	go startServer(portCH)
	port := <-portCH
	startClient(port)
}
