package service

import (
	"block-service/global"
	"block-service/internal/dao"
	"context"
)

type Service struct {
	ctx context.Context
	dao *dao.Dao
}

func New(ctx context.Context) Service {
	return Service{
		ctx,
		dao.New(global.DBEngine),
	}
}


