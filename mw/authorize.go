package mw

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"net/http"
	"slack-clone-api/auth"
	"slack-clone-api/domain/user"
	"strings"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type getTokenFunc func(string, context.Context) (string, error)

func (fn getTokenFunc) GetToken(ID string, ctx context.Context) (string, error) {
	return fn(ID, ctx)
}

func JWTConfig(sign string, svc getTokenFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		authKey := c.Request.Header.Get("Authorization")
		tokenString := strings.TrimPrefix(authKey, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(sign), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		if payload := auth.GetTokenClaims(token); payload != nil {
			at, err := svc.GetToken(payload.UserID, c)
			if err != nil || at != tokenString {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "unauthorized",
				})
				return
			}
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ID := claims["user_id"].(string)
			UUID, _ := uuid.Parse(ID)
			user := &user.User{
				ID:   UUID,
				Role: user.Role(claims["role"].(string)),
			}
			c.Set("user", user)
		}

		c.Next()
	}
}

func ValidatorOnlyAPIKey(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKeyHeader := c.Request.Header.Get("X-API-KEY")
		apiKeyEnc := b64.StdEncoding.EncodeToString([]byte(apiKey))

		if apiKeyHeader != apiKeyEnc {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}

		c.Next()
	}
}
