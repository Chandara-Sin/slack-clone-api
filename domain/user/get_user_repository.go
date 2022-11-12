package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

func GetUser(db *bun.DB) func(uuid.UUID, context.Context) (User, error) {
	return func(ID uuid.UUID, ctx context.Context) (User, error) {
		usr := User{}
		err := db.NewSelect().Model(&usr).Where("id = ?", ID).Scan(ctx)
		return usr, err
	}
}
