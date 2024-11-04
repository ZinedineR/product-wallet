package service

import (
	"context"
	"product-wallet/internal/model"
	"product-wallet/pkg/exception"
)

type WalletService interface {
	// CRUD operations for Wallet
	Create(
		ctx context.Context, req *model.CreateWalletReq,
	) (*model.CreateWalletRes, *exception.Exception)
	Update(
		ctx context.Context, req *model.UpdateWalletReq,
	) (*model.UpdateWalletRes, *exception.Exception)
	DetailWalletTransaction(ctx context.Context, req model.GetWalletByTransactionReq) (
		*model.GetWalletByTransactionRes, *exception.Exception,
	)
	Find(ctx context.Context, req *model.GetAllWalletReq) (*model.GetAllWalletRes, *exception.Exception)
	Detail(ctx context.Context, req *model.GetWalletByIDReq) (*model.GetWalletByIDRes, *exception.Exception)
	Delete(ctx context.Context, req *model.DeleteWalletReq) (*model.DeleteWalletRes, *exception.Exception)
}
