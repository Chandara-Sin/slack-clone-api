package auth

import (
	"context"
	"slack-clone-api/domain/user"
	"time"

	b64 "encoding/base64"

	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"
)

type AuthRepository struct {
	DB  *bun.DB
	RDB *redis.Client
}

func (a AuthRepository) InsertUserByEmail(eml string, ctx context.Context) (user.User, error) {
	usr := user.User{
		FirstName: "",
		LastName:  "",
		Email:     eml,
		Role:      user.Member,
	}
	_, err := a.DB.NewInsert().Model(&usr).Ignore().Exec(ctx)
	return usr, err
}

func (a AuthRepository) SetAuthToken(key string, authCode string, ctx context.Context) (string, error) {
	token := b64.StdEncoding.EncodeToString([]byte(key))
	now := time.Now()
	durat := now.Add(5 * time.Minute)
	err := a.RDB.Set(ctx, token, authCode, durat.Sub(now)).Err()
	return token, err
}

func (a AuthRepository) GetAuthCode(token string, ctx context.Context) (string, error) {
	rs, err := a.RDB.Get(ctx, token).Result()
	return rs, err
}

func (a AuthRepository) ClearAuthCode(key string, ctx context.Context) error {
	rs := a.RDB.Del(ctx, key)
	if rs.Err() == redis.Nil {
		return nil
	}
	return rs.Err()
}
