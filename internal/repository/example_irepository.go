package repository

import (
	"boiler-plate-clean/internal/entity"
)

type ExampleRepository interface {
	// Example operations
	CommonQuery[entity.Example]
}
