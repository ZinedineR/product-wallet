package repository

import (
	"product-wallet/internal/entity"
)

type ExampleRepository interface {
	// Example operations
	CommonQuery[entity.Example]
}
