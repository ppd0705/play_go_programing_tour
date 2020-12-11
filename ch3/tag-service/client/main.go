package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	pb "tag-service/proto"
)

func main() {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "localhost:1034", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("dial err:%v", err)
	}
	defer conn.Close()
	client := pb.NewTagServiceClient(conn)
	req := pb.GetTagListRequest{Name: "Golang"}
	resp, err := client.GetTagList(ctx, &req)
	if err != nil {
		log.Fatalf("GetTagList err: %v", err)
	}
	log.Printf("resp: %v", resp)
}
