package auth

import (
	"context"
	"fmt"
	"math"
	"slack-clone-api/domain/user"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type AuthStore struct {
	DB  *gorm.DB
	RDB *redis.Client
}

func (a AuthStore) GetUser(ID string) (user.User, error) {
	usr := user.User{}
	r := a.DB.First(&usr, ID)
	return usr, r.Error
}

func (a AuthStore) GetUserByEmail(eml string) (user.User, error) {
	usr := user.User{}
	err := a.DB.Where("email = ?", eml).First(&usr).Error
	return usr, err
}

func (a AuthStore) SetToken(ID string, token *AuthToken) error {
	now := time.Now()
	at, _ := ValidateToken(token.AccessToken)
	claims := GetTokenClaims(at)
	atDuration := claims.ExpiresAt.Time
	err := a.RDB.Set(context.TODO(), ID, token.AccessToken, atDuration.Sub(now))
	return err.Err()
}

func (a AuthStore) GetToken(ID string) (string, error) {
	rs, err := a.RDB.Get(context.TODO(), ID).Result()
	return rs, err
}

func GenerateJWTPair(usr user.User) (*AuthToken, error) {
	atClaims := &JwtCustomClaims{
		UserID: strconv.FormatUint(uint64(usr.Id), 10),
		Role:   usr.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "https:slack-clone-api",
			Audience:  jwt.ClaimStrings{"Slack Auth Api"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	ats, err := at.SignedString([]byte(viper.GetString("jwt.secret")))
	if err != nil {
		return nil, err
	}

	rfClaims := &JwtCustomClaims{
		UserID: strconv.FormatUint(uint64(usr.Id), 10),
		Role:   usr.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "https:slack-clone-api",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	rf := jwt.NewWithClaims(jwt.SigningMethodHS256, rfClaims)
	rfs, err := rf.SignedString([]byte(viper.GetString("jwt.secret")))
	if err != nil {
		return nil, err
	}

	return &AuthToken{
		AccessToken:  ats,
		RefreshToken: rfs,
	}, nil
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
		payload.ID = claims["id"].(string)
		payload.Role = user.Role(claims["role"].(string))

		integ, decim := math.Modf(claims["exp"].(float64))
		time := time.Unix(int64(integ), int64(decim*(1e9)))
		payload.ExpiresAt = jwt.NewNumericDate(time)
	}
	return &payload
}
