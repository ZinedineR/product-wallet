package entity

import (
	"os"
	"time"
)

const (
	WalletTableName = "wallet"
)

type Wallet struct {
	Id              string     `json:"id" gorm:"primaryKey;type:uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name            string     `json:"name" example:"personal"`
	UserId          string     `bson:"user_id" json:"user_id" validate:"required,uuid" gorm:"type:uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	User            *User      `bson:"user" json:"user" gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Balance         float64    `gorm:"default:0" json:"balance"`
	LastTransaction *time.Time `gorm:"autoUpdateTime" json:"last_transaction"`
}

func (model *Wallet) Increase(amount float64) {
	model.Balance += amount
}

func (model *Wallet) Decrease(amount float64) {
	model.Balance -= amount
}

func (model *Wallet) TableName() string {
	return os.Getenv("DB_PREFIX") + WalletTableName
}
