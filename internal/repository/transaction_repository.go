package repository

import (
	"product-wallet/internal/entity"
)

type TransactionSQLRepo struct {
	Repository[entity.Transaction]
}

func NewTransactionSQLRepository() TransactionRepository {
	return &TransactionSQLRepo{}
}
