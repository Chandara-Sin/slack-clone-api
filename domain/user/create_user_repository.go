package user

import (
	"gorm.io/gorm"
)

func Create(db *gorm.DB) func(User) error {
	return func(usr User) error {
		r := db.Create(&usr)
		return r.Error
	}
}
