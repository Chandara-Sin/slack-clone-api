package auth

import (
	"slack-clone-api/domain/user"

	"github.com/golang-jwt/jwt/v5"
)

type GrantType string

const (
	Code       GrantType = "auth_code"
	VerifyCode GrantType = "verify_code"
)

type SignUp struct {
	Email     string    `json:"email" binding:"required"`
	GrantType GrantType `json:"grant_type" binding:"required"`
}

type AuthCode struct {
	Code      string    `json:"auth_code" binding:"required"`
	Token     string    `json:"token" binding:"required"`
	GrantType GrantType `json:"grant_type" binding:"required"`
}

type JwtCustomClaims struct {
	UserID string    `json:"user_id"`
	Role   user.Role `json:"role"`
	Email  string    `json:"email"`
	Name   string    `json:"name"`
	jwt.RegisteredClaims
}
