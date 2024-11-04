package model

import (
	"github.com/google/uuid"
	"product-wallet/internal/entity"
)

type BaseProductReq struct {
	Name        string  `json:"name" validate:"required"`
	Price       float64 `json:"price" validate:"required"`
	Description string  `json:"description"`
	Quantity    uint    `json:"quantity" validate:"required"`
	Available   bool    `json:"available"`
}

type CreateProductReq struct {
	BaseProductReq
}

func (req BaseProductReq) ToEntity() *entity.Product {

	if req.Quantity == 0 {
		req.Available = false
	} else {
		req.Available = true
	}
	return &entity.Product{
		Id:          uuid.NewString(),
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
		Quantity:    req.Quantity,
		Available:   req.Available,
	}
}

type CreateProductRes struct {
	entity.Product
}

type UpdateProductReq struct {
	BaseProductReq
	ID string `swaggerignore:"true"`
}
type UpdateProductRes struct {
	entity.Product
}

type DeleteProductReq struct {
	ID string `swaggerignore:"true"`
}
type DeleteProductRes struct {
	ID string `swaggerignore:"true"`
}

type GetAllProductReq struct {
	Page   PaginationParam
	Filter FilterParams
	Sort   OrderParam
}
type GetAllProductRes struct {
	PaginationData[entity.Product]
}

type GetProductByIDReq struct {
	ID string `swaggerignore:"true"`
}

type GetProductByIDRes struct {
	entity.Product
}
