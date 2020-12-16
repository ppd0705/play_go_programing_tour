package main

import (
	"context"
	"flag"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"strings"
	"tag-service/internal/middleware"
	pb "tag-service/proto"
	"tag-service/server"
)

const ServiceName = "tag-service"

var port string

func init() {
	flag.StringVar(&port, "port", "1034", "listen port")
	flag.Parse()
}
func grpcHandleerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}

func runHttpServer() *http.ServeMux {
	serverMux := http.NewServeMux()
	serverMux.Handle("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	}))
	return serverMux
}

func runGrpcGatewayServer() *runtime.ServeMux {
	endpoint := "0.0.0.0:" + port
	gwmux := runtime.NewServeMux()
	dopts := []grpc.DialOption{grpc.WithInsecure()}
	_ = pb.RegisterTagServiceHandlerFromEndpoint(context.Background(), gwmux, endpoint, dopts)
	return gwmux
}

func runGrpcServer() *grpc.Server {
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			middleware.Recovery,
			HelloInterceptor,
			middleware.AccessLog,
			middleware.ErrorLog,
			)),
	}
	s := grpc.NewServer(opts...)
	pb.RegisterTagServiceServer(s, server.NewTagServer())
	reflection.Register(s)
	return s
}

func RunServer(port string) error {
	httpMux := runHttpServer()
	gatewayMux := runGrpcGatewayServer()
	grpcS := runGrpcServer()
	httpMux.Handle("/", gatewayMux)

	//etcdClient, err := clientv3.New(clientv3.Config{
	//	Endpoints: []string{"http://localhost:2379"},
	//	DialTimeout: time.Second*3,
	//})
	//if err != nil {
	//	return err
	//}
	//defer etcdClient.Close()
	//target := fmt.Sprintf("/etcdv3://golang/grpc/%s", ServiceName)
	//grpcproxy.Register(etcdClient, target, ":"+port, 60)

	return http.ListenAndServe(":"+port, grpcHandleerFunc(grpcS, httpMux))

}
func HelloInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println("in")
	resp, err := handler(ctx, req)
	log.Println("exit")
	return resp, err
}

func main() {
	err := RunServer(port)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
}
func runServer0() {
	s := grpc.NewServer()
	pb.RegisterTagServiceServer(s, server.NewTagServer())
	reflection.Register(s)


	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}
	fmt.Printf("start listen: %s\n", lis.Addr().String())
	if err = s.Serve(lis); err != nil {
		log.Fatalf("server.Serve err: %v", err)
	}

}
