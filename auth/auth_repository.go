package auth

import (
	"context"
	"fmt"
	"math"
	"slack-clone-api/domain/user"
	"time"

	"github.com/go-redis/redis/v8"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
)

type AuthRepository struct {
	DB  *bun.DB
	RDB *redis.Client
}

func (a AuthRepository) GetUser(ID string, ctx context.Context) (user.User, error) {
	usr := user.User{}
	err := a.DB.NewSelect().Model(&usr).Where("id = ?", ID).Scan(ctx)
	return usr, err
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

func (a AuthRepository) InsertAuthToken(authCode string, ctx context.Context) (string, error) {
	token := uuid.New().String()
	now := time.Now()
	durat := now.Add(5 * time.Minute)
	err := a.RDB.Set(ctx, token, authCode, durat.Sub(now)).Err()
	return token, err

}

func (a AuthRepository) SetAuthToken(ID string, token *AuthToken, ctx context.Context) error {
	now := time.Now()
	if err := setToken(ctx, a.RDB, ID, token.AccessToken, now); err != nil {
		return err
	}
	rf, _ := ValidateToken(token.RefreshToken)
	rfClaims := GetTokenClaims(rf)
	err := setToken(ctx, a.RDB, rfClaims.Subject, token.RefreshToken, now)
	return err
}

func (a AuthRepository) GetToken(ID string, ctx context.Context) (string, error) {
	rs, err := a.RDB.Get(ctx, ID).Result()
	return rs, err
}

func (a AuthRepository) ClearToken(key string, ctx context.Context) error {
	rs := a.RDB.Del(ctx, key)
	return rs.Err()
}

func GenerateJWTPair(usr user.User) (*AuthToken, error) {
	atJti := uuid.New().String()
	rfJti := uuid.New().String()

	at, err := generateToken(usr, atJti, rfJti, (5 * time.Minute))
	if err != nil {
		return nil, err
	}
	rf, err := generateToken(usr, rfJti, atJti, (10 * time.Minute))
	if err != nil {
		return nil, err
	}

	return &AuthToken{
		AccessToken:  at,
		RefreshToken: rf,
	}, nil
}

func generateToken(usr user.User, ID string, subj string, durat time.Duration) (string, error) {
	claims := &JwtCustomClaims{
		UserID: usr.ID.String(),
		Role:   usr.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        ID,
			Subject:   subj,
			Issuer:    "https://slack-clone-api",
			Audience:  jwt.ClaimStrings{"Slack Auth Api"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(durat)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tk, err := jwtToken.SignedString([]byte(viper.GetString("jwt.secret")))
	if err != nil {
		return "", err
	}
	return tk, nil
}

func setToken(ctx context.Context, rdb *redis.Client, key string, token string, exp time.Time) error {
	jwtToken, _ := ValidateToken(token)
	claims := GetTokenClaims(jwtToken)
	durat := claims.ExpiresAt.Time
	return rdb.Set(ctx, key, token, durat.Sub(exp)).Err()
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(viper.GetString("jwt.secret")), nil
	})
	return token, err
}

func GetTokenClaims(token *jwt.Token) *JwtCustomClaims {
	payload := JwtCustomClaims{}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		payload.UserID = claims["user_id"].(string)
		payload.Role = user.Role(claims["role"].(string))
		payload.ID = claims["jti"].(string)
		payload.Subject = claims["sub"].(string)

		integ, decim := math.Modf(claims["exp"].(float64))
		time := time.Unix(int64(integ), int64(decim*(1e9)))
		payload.ExpiresAt = jwt.NewNumericDate(time)
	}
	return &payload
}
