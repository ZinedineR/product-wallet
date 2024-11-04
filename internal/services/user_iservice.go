package service

import (
	"context"
	"product-wallet/internal/model"
	"product-wallet/pkg/exception"
)

type UserService interface {
	// CRUD operations for User
	Register(
		ctx context.Context, req *model.CreateUserReq,
	) (*model.CreateUserRes, *exception.Exception)
	Login(ctx context.Context, req *model.CreateUserReq) (*model.LoginUserRes, *exception.Exception)
}
