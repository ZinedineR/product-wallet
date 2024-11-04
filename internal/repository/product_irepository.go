package repository

import (
	"product-wallet/internal/entity"
)

type ProductRepository interface {
	CommonQuery[entity.Product]
}
