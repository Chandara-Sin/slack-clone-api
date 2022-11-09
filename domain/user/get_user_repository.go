package user

import (
	"context"

	"github.com/uptrace/bun"
)

func GetUser(db *bun.DB) func(string) (User, error) {
	return func(ID string) (User, error) {
		usr := User{}
		err := db.NewSelect().Model(&usr).Where("id = ?", ID).Scan(context.TODO())
		return usr, err
	}
}
