package repository

import (
	"product-wallet/internal/entity"
)

type WalletRepository interface {
	CommonQuery[entity.Wallet]
}
