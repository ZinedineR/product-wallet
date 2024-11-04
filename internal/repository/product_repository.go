package repository

import (
	"product-wallet/internal/entity"
)

type ProductSQLRepo struct {
	Repository[entity.Product]
}

func NewProductSQLRepository() ProductRepository {
	return &ProductSQLRepo{}
}
