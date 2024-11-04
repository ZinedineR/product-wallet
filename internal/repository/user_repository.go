package repository

import (
	"product-wallet/internal/entity"
)

type UserSQLRepo struct {
	Repository[entity.User]
}

func NewUserSQLRepository() UserRepository {
	return &UserSQLRepo{}
}
