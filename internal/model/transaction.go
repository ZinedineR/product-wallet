package model

import (
	"github.com/google/uuid"
	"product-wallet/internal/entity"
	"product-wallet/pkg/utils/converter"
)

type BaseTransactionReq struct {
	WalletId        string  `json:"wallet_id" validate:"required,uuid"`
	ProductId       *string `json:"product_id,omitempty" validate:"required,uuid"`
	ProductQuantity *uint   `json:"product_quantity,omitempty" validate:"required,number"`
}

type CreateTransactionReq struct {
	BaseTransactionReq
}

func (req BaseTransactionReq) ToEntity() *entity.Transaction {
	return &entity.Transaction{
		Id:        uuid.NewString(),
		ProductId: req.ProductId,
		WalletId:  req.WalletId,
		Type:      "expense",
	}
}

//func (req BaseTransactionReq) ToEntity(productname string, totalprice float64) *entity.Transaction {
//	if req.ProductQuantity == nil {
//		var one uint = 1
//		req.ProductQuantity = &one
//	}
//	return &entity.Transaction{
//		Id:          uuid.NewString(),
//		ProductId:   req.ProductId,
//		WalletId:    req.WalletId,
//		Type:        "expense",
//		Description: "Buying " + productname + ", quantity: " + converter.ToString(*req.ProductQuantity) + " for " + converter.ToString(totalprice),
//		Amount:      totalprice,
//	}
//}

func (req BaseTransactionReq) ToProductEntity(product *entity.Product) *entity.Product {
	var (
		quantity uint
	)
	quantity = product.Quantity - *req.ProductQuantity
	if quantity == 0 {
		product.Available = false
	} else {
		product.Available = true
	}
	return &entity.Product{
		Id:          product.Id,
		Name:        product.Name,
		Price:       product.Price,
		Description: product.Description,
		Quantity:    product.Quantity - *req.ProductQuantity,
		Available:   product.Available,
	}
}

type CreateTransactionRes struct {
	entity.Transaction
}

type UpdateTransactionReq struct {
	BaseTransactionReq
	ID string
}
type UpdateTransactionRes struct {
	entity.Transaction
}

type DeleteTransactionReq struct {
	ID string `swaggerignore:"true"`
}
type DeleteTransactionRes struct {
	ID string `swaggerignore:"true"`
}

type GetAllTransactionReq struct {
	Page   PaginationParam
	Filter FilterParams
	Sort   OrderParam
}
type GetAllTransactionRes struct {
	PaginationData[entity.Transaction]
}

type GetTransactionByIDReq struct {
	ID string `swaggerignore:"true"`
}

type GetTransactionByIDRes struct {
	entity.Transaction
}

type CreditTransactionReq struct {
	WalletId string  `json:"wallet_id" validate:"required"`
	Amount   float64 `json:"amount" validate:"required"`
}
type CreditTransactionRes struct {
	entity.Transaction
}

func (req CreditTransactionReq) ToEntity() *entity.Transaction {
	return &entity.Transaction{
		Id:          uuid.NewString(),
		WalletId:    req.WalletId,
		Type:        "income",
		Amount:      req.Amount,
		Description: "Credit of " + converter.ToString(req.Amount),
	}
}

type TransferTransactionReq struct {
	SenderId   string  `json:"wallet_id" validate:"required"`
	ReceiverId string  `json:"receiver_id" validate:"required"`
	Amount     float64 `json:"amount" validate:"required"`
}
type TransferTransactionRes struct {
	SenderTransaction   entity.Transaction `json:"sender_transaction"`
	ReceiverTransaction entity.Transaction `json:"receiver_transaction"`
}

func (req TransferTransactionReq) ToSenderEntity(receiverName, senderWalletID string) *entity.Transaction {
	return &entity.Transaction{
		Id:          uuid.NewString(),
		Type:        "transfer",
		Amount:      req.Amount,
		Description: "Transfer to: " + receiverName,
		WalletId:    senderWalletID,
	}
}

func (req TransferTransactionReq) ToReceiverEntity(senderName, receiverWalletID string) *entity.Transaction {
	return &entity.Transaction{
		Id:          uuid.NewString(),
		Type:        "transfer",
		Amount:      req.Amount,
		Description: "Transfer from: " + senderName,
		WalletId:    receiverWalletID,
	}
}
