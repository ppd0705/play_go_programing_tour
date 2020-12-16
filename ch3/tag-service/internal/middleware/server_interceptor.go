package middleware

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"runtime/debug"
	"tag-service/pkg/errcode"
	"time"
)

func AccessLog(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	requestLog := "access request: method: %s, begin_time: %d, request: %v"
	beginTime := time.Now().Local().Unix()

	log.Printf(requestLog, info.FullMethod, beginTime, req)
	resp, err := handler(ctx, req)
	responseLog := "access response: method: %s, begin_time: %d, end_time: %d, response: %v"
	endTime := time.Now().Local().Unix()
	log.Printf(responseLog, info.FullMethod, beginTime, endTime, resp)
	return resp, err
}

func ErrorLog(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	resp, err := handler(ctx, req)
	if err != nil {
		errLog := "error log: method: %s, code: %v, message: %v, details %v"
		s := errcode.FromError(err)
		log.Printf(errLog, info.FullMethod, s.Code(), s.Message(), s.Details())
	}
	return resp, err
}

func Recovery(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	defer func() {
		if err := recover(); err != nil {
			recoverLog := "recovery log: method: %s, message: %v, stack: %s"
			log.Printf(recoverLog, info.FullMethod, err, string(debug.Stack()))
		}
	}()
	return handler(ctx, req)
}