package main

import (
	"context"
	"google.golang.org/grpc"
	pb "grpc-demo/proto"
	"grpc-demo/server"
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

func startClient(port string) {
	conn, _ := grpc.Dial(port, grpc.WithInsecure())
	defer conn.Close()

	client := pb.NewGreeterClient(conn)
	err := SayHello(client)
	if err != nil {
		log.Fatalf("Say hello err: %v", err)
	}
}
func main() {
	portCH := make(chan string)
	go startServer(portCH)
	port := <-portCH
	startClient(port)
}
