package service

import (
	"context"
	"product-wallet/internal/model"
	"product-wallet/pkg/exception"
)

type ProductService interface {
	// CRUD operations for Product
	Create(
		ctx context.Context, req *model.CreateProductReq,
	) (*model.CreateProductRes, *exception.Exception)
	Update(
		ctx context.Context, req *model.UpdateProductReq,
	) (*model.UpdateProductRes, *exception.Exception)
	Find(ctx context.Context, req *model.GetAllProductReq) (*model.GetAllProductRes, *exception.Exception)
	Detail(ctx context.Context, req *model.GetProductByIDReq) (*model.GetProductByIDRes, *exception.Exception)
	Delete(ctx context.Context, req *model.DeleteProductReq) (*model.DeleteProductRes, *exception.Exception)
}
