package auth

import (
	"slack-clone-api/domain/user"

	"github.com/golang-jwt/jwt/v4"
)

type GrantType string

const (
	AuthCode   GrantType = "auth_code"
	VerifyCode GrantType = "verify_code"
)

type Login struct {
	Email     string    `json:"email,omitempty"`
	AuthCode  string    `json:"auth_code,omitempty"`
	GrantType GrantType `json:"grant_type" binding:"required"`
}

type SignOut struct {
	Token         string `json:"token"`
	TokenTypeHint string `json:"token_type_hint"`
}

type JwtCustomClaims struct {
	UserID string    `json:"user_id"`
	Role   user.Role `json:"role"`
	jwt.RegisteredClaims
}

type AuthToken struct {
	AccessToken  string
	RefreshToken string
}
