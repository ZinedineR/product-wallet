package service

import (
	"context"
	"product-wallet/internal/model"
	"product-wallet/pkg/exception"
)

type TransactionService interface {
	// CRUD operations for Example
	Create(
		ctx context.Context, req *model.CreateTransactionReq,
	) (*model.CreateTransactionRes, *exception.Exception)
	Detail(ctx context.Context, req *model.GetTransactionByIDReq) (*model.GetTransactionByIDRes, *exception.Exception)
	Find(ctx context.Context, req *model.GetAllTransactionReq) (
		*model.GetAllTransactionRes, *exception.Exception,
	)
	Credit(
		ctx context.Context, req *model.CreditTransactionReq,
	) (*model.CreditTransactionRes, *exception.Exception)
	Transfer(
		ctx context.Context, req *model.TransferTransactionReq,
	) (*model.TransferTransactionRes, *exception.Exception)
	Delete(ctx context.Context, req *model.DeleteTransactionReq) (*model.DeleteTransactionRes, *exception.Exception)
}
