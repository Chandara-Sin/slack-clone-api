package auth

import (
	"context"
	"errors"
	"net/http"
	"slack-clone-api/domain/user"
	"slack-clone-api/logger"

	"github.com/gin-gonic/gin"

	"golang.org/x/crypto/bcrypt"
)

func JWTConfigHandler(svc AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.Unwrap(c)

		reqLogin := Login{}
		if err := c.ShouldBindJSON(&reqLogin); err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		usr := user.User{}
		if reqLogin.GrantType == Password {
			rs, err := validateUser(reqLogin, svc, c)
			if err != nil {
				log.Error(err.Error())
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "username or password is incorrect",
				})
				return
			}

			if token, _ := svc.GetToken(rs.ID.String(), c); token != "" {
				at, _ := ValidateToken(token)
				atClaims := GetTokenClaims(at)
				svc.ClearToken(atClaims.UserID, c)
				svc.ClearToken(atClaims.ID, c)
			}

			usr = rs
		} else if reqLogin.GrantType == RefreshToken {
			token, err := ValidateToken(reqLogin.RefreshToken)
			if err != nil {
				log.Error(err.Error())
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": err.Error(),
				})
				return
			}

			claims := GetTokenClaims(token)

			if _, err := svc.GetToken(claims.Subject, c); err != nil {
				log.Error(err.Error())
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "unauthorized",
				})
				return
			} else {
				svc.ClearToken(claims.UserID, c)
				svc.ClearToken(claims.Subject, c)
			}

			res, err := svc.GetUser(claims.UserID, c)
			if err != nil {
				log.Error(err.Error())
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": err.Error(),
				})
				return
			}
			usr = res
		}

		authToken, err := GenerateJWTPair(usr)
		if err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		if svc.SetToken(usr.ID.String(), authToken, c); err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token":  authToken.AccessToken,
			"refresh_token": authToken.RefreshToken,
			"token_type":    "Bearer",
		})

	}
}

func validateUser(reqLogin Login, svc AuthService, ctx context.Context) (user.User, error) {
	usr, err := svc.GetUserByEmail(reqLogin.Email, ctx)
	if err != nil {
		return usr, err
	}

	match := checkPasswordHash(reqLogin.Password, usr.HashedPassword)
	if !match {
		return usr, errors.New("incorrect password")
	}
	return usr, nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
