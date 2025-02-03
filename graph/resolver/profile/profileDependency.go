package profile

import (
	"gorm.io/gorm"
)

type ProfileResolver struct {
	db *gorm.DB
}

func InitializeProfileResolver( db *gorm.DB) *ProfileResolver {
	return &ProfileResolver{
		db: db,
	}
}

const ProfilesKey = "profiles"
const ProfileKey = "profile:"
