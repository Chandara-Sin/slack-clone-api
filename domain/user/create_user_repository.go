package user

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

func Create(db *bun.DB) func(User) error {
	return func(usr User) error {
		r, err := db.NewInsert().Model(&usr).Exec(context.TODO())
		fmt.Println(r)
		return err
	}
}
