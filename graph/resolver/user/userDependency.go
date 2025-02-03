package user

import (
	"gorm.io/gorm"
)

type UserResolver struct{
	db *gorm.DB
}

func InitializeUserResolver( db *gorm.DB) *UserResolver {
	return &UserResolver{
		db: db,
	}
}