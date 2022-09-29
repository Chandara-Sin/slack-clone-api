package auth

import (
	"slack-clone-api/domain/user"

	"github.com/golang-jwt/jwt/v4"
)

type GrantType string

const (
	Password     GrantType = "password"
	RefreshToken GrantType = "refresh_token"
)

type Login struct {
	Email        string    `json:"email,omitempty"`
	Password     string    `json:"password,omitempty"`
	GrantType    GrantType `json:"grant_type" binding:"required"`
	RefreshToken string    `json:"refresh_token,omitempty"`
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
	GetUserByEmail(eml string) (user.User, error)
	GetUser(ID string) (user.User, error)
	SetToken(ID string, token *AuthToken) error
	GetToken(ID string) (string, error)
}
