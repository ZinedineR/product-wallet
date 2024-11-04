package model

import (
	"github.com/google/uuid"
	"product-wallet/internal/entity"
)

type BaseUserReq struct {
	Username string `json:"username" validate:"required" example:"john_doe"`
	Password string `json:"password" validate:"required,password,gte=8" example:"SecurePass123!"`
}

type CreateUserReq struct {
	BaseUserReq
}

func (req BaseUserReq) ToEntity(password string) *entity.User {
	return &entity.User{
		Id:       uuid.NewString(),
		Username: req.Username,
		Password: password,
	}
}

type CreateUserRes struct {
	entity.User
}

type LoginUserRes struct {
	Username string `json:"username" example:"john_doe"`
	Token    string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"` // JWT token example

}

type UpdateUserReq struct {
	BaseUserReq
	ID string
}
type UpdateUserRes struct {
	entity.User
}

type DeleteUserReq struct {
	BaseUserReq
	ID string `json:"id" name:"id"`
}
type DeleteUserRes struct {
	ID string `json:"id" name:"id"`
}

type GetAllUserReq struct {
	Page   PaginationParam
	Filter FilterParams
	Sort   OrderParam
}
type GetAllUserRes struct {
	PaginationData[entity.User]
}

type GetUserByIDReq struct {
	ID string `json:"id" name:"id"`
}

type GetUserByIDRes struct {
	entity.User
}
