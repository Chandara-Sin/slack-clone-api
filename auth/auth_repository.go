package auth

import (
	"context"
	"slack-clone-api/domain/user"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
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

func (a AuthRepository) SetAuthToken(authCode string, ctx context.Context) (string, error) {
	token := uuid.New().String()
	now := time.Now()
	durat := now.Add(5 * time.Minute)
	err := a.RDB.Set(ctx, token, authCode, durat.Sub(now)).Err()
	return token, err

}

func (a AuthRepository) GetAuthCode(ID string, ctx context.Context) (string, error) {
	rs, err := a.RDB.Get(ctx, ID).Result()
	return rs, err
}

func (a AuthRepository) ClearAuthCode(key string, ctx context.Context) error {
	rs := a.RDB.Del(ctx, key)
	return rs.Err()
}
