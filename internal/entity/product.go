package entity

import (
	"os"
	"time"
)

const (
	ProductTableName = "product"
)

type Product struct {
	Id          string     `json:"id" gorm:"primaryKey;type:uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name        string     `json:"name"`
	Price       float64    `json:"price"`
	Description string     `json:"description"`
	Quantity    uint       `json:"quantity"`
	Available   bool       `json:"available"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

func (model *Product) TableName() string {
	return os.Getenv("DB_PREFIX") + ProductTableName
}
