package main

import (
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"log"
	"tag-service/internal/middleware"
	pb "tag-service/proto"
)

type Auth struct {
	AppKey string
	APpSecret string
}

func (a *Auth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error){
	return map[string]string{"app_key": a.AppKey, "app_secret": a.APpSecret}, nil
}

func (a *Auth) RequireTransportSecurity() bool {
	return false
}

func GetClientConn(ctx context.Context, serviceName string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	//config := clientv3.Config{
	//	Endpoints: []string{"http://localhost:2379"},
	//	DialTimeout: time.Second * 60,
	//}
	//cli, err := clientv3.New(config)
	//if err != nil {
	//	return nil, err
	//}
	//r := &naming.GRPCResolver{Client: cli}
	target := fmt.Sprintf("/etcdv3://golang/grpc/%s", serviceName)
	//opts = append(opts, grpc.WithBalancer(grpc.RoundRobin(r)), grpc.WithBlock())
	return grpc.DialContext(ctx, target, opts...)
}


func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure(),grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			grpc_retry.UnaryClientInterceptor(grpc_retry.WithMax(2), grpc_retry.WithCodes(codes.Unknown, codes.Internal, codes.DeadlineExceeded)),
			middleware.UnaryContextTimeout(),
			)))

	auth := Auth{AppKey: "tag server", APpSecret: "hello world"}

	opts = append(opts, grpc.WithPerRPCCredentials(&auth))
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "localhost:1034",  opts...)
	//conn, err := GetClientConn(ctx, "tag-service", opts)
	if err != nil {
		log.Fatalf("dial err:%v", err)
	}
	defer conn.Close()
	client := pb.NewTagServiceClient(conn)
	req := pb.GetTagListRequest{Name: "Golang"}
	newCtx := metadata.AppendToOutgoingContext(ctx, "hello", "world")
	resp, err := client.GetTagList(newCtx, &req)
	if err != nil {
		log.Fatalf("GetTagList err: %v", err)
	}
	log.Printf("resp: %v", resp)
}
