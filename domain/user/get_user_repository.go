package user

import (
	"context"

	"github.com/uptrace/bun"
)

func GetUser(db *bun.DB) func(string, context.Context) (User, error) {
	return func(ID string, ctx context.Context) (User, error) {
		usr := User{}
		err := db.NewSelect().Model(&usr).Where("id = ?", ID).Scan(ctx)
		return usr, err
	}
}
