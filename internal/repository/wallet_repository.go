package repository

import (
	"product-wallet/internal/entity"
)

type WalletSQLRepo struct {
	Repository[entity.Wallet]
}

func NewWalletSQLRepository() WalletRepository {
	return &WalletSQLRepo{}
}
