package model

import (
	"github.com/google/uuid"
	"product-wallet/internal/entity"
	"time"
)

type BaseWalletReq struct {
	Name   string `json:"name" example:"personal"`
	UserId string `json:"user_id" validate:"required,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
}

type CreateWalletReq struct {
	BaseWalletReq
}

func (req BaseWalletReq) ToEntity() *entity.Wallet {
	return &entity.Wallet{
		Id:     uuid.NewString(),
		Name:   req.Name,
		UserId: req.UserId,
	}
}

type CreateWalletRes struct {
	entity.Wallet
}

type UpdateWalletReq struct {
	BaseWalletReq
	ID string `json:"id" name:"id"`
}
type UpdateWalletRes struct {
	entity.Wallet
}

type DeleteWalletReq struct {
	ID string `json:"id" name:"id"`
}
type DeleteWalletRes struct {
	ID string `json:"id" name:"id"`
}

type GetAllWalletReq struct {
	Page   PaginationParam
	Filter FilterParams
	Sort   OrderParam
}
type GetAllWalletRes struct {
	PaginationData[entity.Wallet]
}

type GetWalletByIDReq struct {
	ID string `json:"id" name:"id"`
}

type GetWalletByIDRes struct {
	entity.Wallet
}

type GetWalletByTransactionReq struct {
	ID   string `json:"id" name:"id"`
	From time.Time
	To   time.Time
}

type GetWalletByTransactionRes struct {
	entity.Wallet
	Transaction []entity.Transaction `json:"transaction"`
}

func NewGetWalletByTransactionRes(wallet entity.Wallet, transaction []entity.Transaction) *GetWalletByTransactionRes {
	return &GetWalletByTransactionRes{Wallet: wallet, Transaction: transaction}
}
