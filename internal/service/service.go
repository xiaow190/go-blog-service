package service

import (
	"context"
	"go-programming-book/blog-service/global"
	"go-programming-book/blog-service/internal/dao"
)

type Service struct {
	ctx context.Context
	dao *dao.Dao
}

func New(ctx context.Context) Service {
	svc := Service{ctx: ctx}
	svc.dao = dao.New(global.DbEngine)
	return svc
}
