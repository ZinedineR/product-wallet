package model

import (
	"github.com/google/uuid"
	"product-wallet/internal/entity"
	"product-wallet/pkg/utils/converter"
)

type BaseTransactionReq struct {
	WalletId string  `json:"wallet_id"`
	Type     string  `json:"type" validate:"eq=income|eq=expense|eq=transfer"`
	Amount   float64 `json:"amount"`
}

type CreateTransactionReq struct {
	BaseTransactionReq
}

func (req BaseTransactionReq) ToEntity() *entity.Transaction {
	return &entity.Transaction{
		Id:       uuid.NewString(),
		WalletId: req.WalletId,
		Type:     req.Type,
		Amount:   req.Amount,
	}
}

type CreateTransactionRes struct {
	entity.Transaction
}

type UpdateTransactionReq struct {
	BaseTransactionReq
	ID string `json:"id" name:"id"`
}
type UpdateTransactionRes struct {
	entity.Transaction
}

type DeleteTransactionReq struct {
	ID string `json:"id" name:"id"`
}
type DeleteTransactionRes struct {
	ID string `json:"id" name:"id"`
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
	ID string `json:"id" name:"id"`
}

type GetTransactionByIDRes struct {
	entity.Transaction
}

type CreditTransactionReq struct {
	WalletId string  `json:"wallet_id"`
	Amount   float64 `json:"amount"`
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
	SenderId   string  `json:"wallet_id"`
	ReceiverId string  `json:"receiver_id"`
	Amount     float64 `json:"amount"`
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
