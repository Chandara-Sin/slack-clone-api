package auth

import (
	"slack-clone-api/domain/user"

	"github.com/golang-jwt/jwt/v4"
)

type JwtCustomClaims struct {
	UserID string    `json:"id"`
	Role   user.Role `json:"role"`
	jwt.RegisteredClaims
}

type AuthToken struct {
	AccessToken  string
	RefreshToken string
}
