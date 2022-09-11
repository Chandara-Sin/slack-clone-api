package auth

import (
	"fmt"
	"slack-clone-api/domain/user"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

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

	rtClaims := &JwtCustomClaims{
		UserID: strconv.FormatUint(uint64(usr.Id), 10),
		Role:   usr.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "https:slack-clone-api",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	rf := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
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
