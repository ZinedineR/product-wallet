package entity

import (
	"os"
)

const (
	UserTableName = "user"
)

type User struct {
	Id       string `json:"id" gorm:"primaryKey;type:uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	Username string `json:"username" example:"john_doe"`
	Password string `json:"password" example:"$2a$12$eixZaYVK1fsbw1ZfbX3OXe.PZyWJQ0Zf10hErsTQ6FVRHiA2vwLHu"` // Example of bcrypt-hashed password
}

func (model *User) TableName() string {
	return os.Getenv("DB_PREFIX") + UserTableName
}
