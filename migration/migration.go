package migration

import (
	"product-wallet/internal/entity"
	"product-wallet/pkg/database"
)

func AutoMigration(CpmDB *database.Database) {
	CpmDB.MigrateDB(
		&entity.User{},
		&entity.Wallet{},
		&entity.Transaction{},
	)
}
