package repository

import (
	"product-wallet/internal/entity"
)

type TransactionRepository interface {
	CommonQuery[entity.Transaction]
}
