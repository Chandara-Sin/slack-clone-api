package user

import (
	"context"

	"github.com/uptrace/bun"
)

func Create(db *bun.DB) func(User, context.Context) error {
	return func(usr User, ctx context.Context) error {
		_, err := db.NewInsert().Model(&usr).Exec(ctx)
		return err
	}
}
