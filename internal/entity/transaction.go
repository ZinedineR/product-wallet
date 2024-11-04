package entity

import (
	"os"
	"time"
)

const (
	TransactionTableName = "transaction"
)

type Transaction struct {
	Id              string     `json:"id" gorm:"primaryKey;type:uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	Type            string     `json:"type" validate:"eq=income|eq=expense|eq=transfer"`
	Amount          float64    `json:"amount"`
	Description     string     `json:"description"`
	WalletId        string     `json:"wallet_id"`
	Wallet          *Wallet    `gorm:"foreignKey:WalletId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"wallet,omitempty"`
	TransactionTime *time.Time `gorm:"autoCreateTime" json:"transaction_time"`
}

func (model *Transaction) TableName() string {
	return os.Getenv("DB_PREFIX") + TransactionTableName
}
