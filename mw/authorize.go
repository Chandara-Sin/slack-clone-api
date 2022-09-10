package mw

import (
	"fmt"
	"net/http"
	"slack-clone-api/auth"
	"slack-clone-api/domain/user"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
)

func JWTConfig(sign string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.Request.Header.Get("Authorization")
		tokenString := strings.TrimPrefix(authorization, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(sign), nil
		})

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(*auth.JwtCustomClaims); ok && token.Valid {
			ID, _ := strconv.ParseUint(claims.UserID, 10, 64)
			user := user.User{
				Id:   uint(ID),
				Role: claims.Role,
			}
			c.Set("user", user)
		}

		c.Next()
	}
}
