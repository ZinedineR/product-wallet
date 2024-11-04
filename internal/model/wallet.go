package model

import (
	"github.com/google/uuid"
	"product-wallet/internal/entity"
	"time"
)

type BaseWalletReq struct {
	Name   string `json:"name" example:"personal"`
	UserId string `json:"-" validate:"required,uuid" swaggerignore:"true"`
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
	ID string `swaggerignore:"true"`
}
type UpdateWalletRes struct {
	entity.Wallet
}

type DeleteWalletReq struct {
	ID string `swaggerignore:"true"`
}
type DeleteWalletRes struct {
	ID string `swaggerignore:"true"`
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
	ID string `swaggerignore:"true"`
}

type GetWalletByIDRes struct {
	entity.Wallet
}

type GetWalletByTransactionReq struct {
	ID   string    `swaggerignore:"true"`
	From time.Time `swaggerignore:"true"`
	To   time.Time `swaggerignore:"true"`
}

type GetWalletByTransactionRes struct {
	entity.Wallet
	Transaction []entity.Transaction `json:"transaction"`
}

func NewGetWalletByTransactionRes(wallet entity.Wallet, transaction []entity.Transaction) *GetWalletByTransactionRes {
	return &GetWalletByTransactionRes{Wallet: wallet, Transaction: transaction}
}
