package auth

import (
	"context"
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

type AuthService interface {
	GetUserByEmail(eml string, ctx context.Context) (user.User, error)
	GetUser(ID string, ctx context.Context) (user.User, error)
	SetAuthToken(ID string, token *AuthToken, ctx context.Context) error
	GetToken(ID string, ctx context.Context) (string, error)
	ClearToken(ID string, ctx context.Context) error
}
