package user

import (
	"context"

	"github.com/uptrace/bun"
)

func Create(db *bun.DB) func(User) error {
	return func(usr User) error {
		_, err := db.NewInsert().Model(&usr).Exec(context.TODO())
		return err
	}
}
