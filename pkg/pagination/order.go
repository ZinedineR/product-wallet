package pagination

import (
	"fmt"
	"product-wallet/internal/model"

	"gorm.io/gorm"
)

func Order(param model.OrderParam, query *gorm.DB) *gorm.DB {
	if param.Order != "" && param.OrderBy != "" {
		query = query.Order(fmt.Sprintf("%s %s", param.OrderBy, param.Order))
	}
	return query
}
