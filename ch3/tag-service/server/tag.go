package server

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/metadata"
	"tag-service/pkg/bapi"
	"tag-service/pkg/errcode"
	pb "tag-service/proto"
)

type TagService struct {
	auth *Auth
}

type Auth struct {
}

func (a *Auth) GetAppKey() string {
	return "tag server"
}

func (a *Auth) GetSecret() string {
	return "hello world"
}

func (a *Auth) Check(ctx context.Context) error {
	md, _ := metadata.FromIncomingContext(ctx)
	var appKey, appSecret string
	if value, ok := md["app_key"]; ok {
		appKey = value[0]
	}
	if value, ok := md["app_secret"]; ok {
		appSecret = value[0]
	}
	if appKey != a.GetAppKey() || appSecret != a.GetSecret() {
		return errcode.TogRPCError(errcode.Unauthorized)
	}
	return nil
}

func NewTagServer() *TagService {
	return &TagService{}
}

var _ pb.TagServiceServer = (*TagService)(nil)

func (t *TagService) GetTagList(ctx context.Context, r *pb.GetTagListRequest) (*pb.GetTagListReply, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	fmt.Printf("md: %v\n", md)
	if err := t.auth.Check(ctx); err != nil {
		return nil, err
	}
	api := bapi.NewAPI("http://localhost:8000")
	body, err := api.GetTagList(ctx, r.GetName())
	if err != nil {
		return nil, errcode.TogRPCError(errcode.ErrorGetTagListFail)
	}
	tagList := pb.GetTagListReply{}
	err = json.Unmarshal(body, &tagList)
	if err != nil {
		return nil, errcode.TogRPCError(errcode.Fail)
	}
	return &tagList, nil
}
