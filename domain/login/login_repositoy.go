package login

import (
	"slack-clone-api/domain/user"

	"gorm.io/gorm"
)

func GetUserByEmail(db *gorm.DB) func(string) (user.User, error) {
	return func(eml string) (user.User, error) {
		usr := user.User{}
		err := db.Where("email = ?", eml).First(&usr).Error
		return usr, err
	}
}
