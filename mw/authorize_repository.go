package mw

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func GetToken(db *redis.Client) func(string, context.Context) (string, error) {
	return func(ID string, ctx context.Context) (string, error) {
		rs, err := db.Get(ctx, ID).Result()
		return rs, err
	}
}
