package repository

import (
	"product-wallet/internal/entity"
)

type UserRepository interface {
	CommonQuery[entity.User]
}
