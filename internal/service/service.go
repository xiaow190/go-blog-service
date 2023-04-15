package service

import (
	"context"
	"go-programming-book/blog-service/global"
	"go-programming-book/blog-service/internal/dao"

	otgorm "github.com/eddycjy/opentracing-gorm"
)

type Service struct {
	ctx context.Context
	dao *dao.Dao
}

func New(ctx context.Context) Service {
	svc := Service{ctx: ctx}
	svc.dao = dao.New(otgorm.WithContext(svc.ctx, global.DbEngine))
	return svc
}
