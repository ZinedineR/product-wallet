package migration

import (
	"boiler-plate-clean/internal/entity"
	"boiler-plate-clean/pkg/database"
)

func AutoMigration(CpmDB *database.Database) {
	CpmDB.MigrateDB(

		&entity.Example{})
	//&entity.SMSLog{}
}
