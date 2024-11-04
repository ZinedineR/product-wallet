package service

import (
	"context"
	"product-wallet/internal/entity"
	"product-wallet/internal/model"
	"product-wallet/pkg/exception"
)

type ExampleService interface {
	// CRUD operations for Example
	CreateExample(
		ctx context.Context, model *entity.Example,
	) *exception.Exception
}

type ListExampleResp struct {
	Pagination *model.Pagination `json:"pagination"`
	Data       []*entity.Example `json:"data"`
}
