package repository

import (
	"product-wallet/internal/entity"
)

type ExampleSQLRepo struct {
	Repository[entity.Example]
}

func NewExampleSQLRepository() ExampleRepository {
	return &ExampleSQLRepo{}
}
