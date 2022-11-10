package auth

import (
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
			res, err := svc.GetUserByEmail(reqLogin.Email)
			if err != nil {
				log.Error(err.Error())
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "username or password is incorrect",
				})
				return
			}

			match := checkPasswordHash(reqLogin.Password, res.HashedPassword)
			if !match {
				log.Error(err.Error())
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "username or password is incorrect",
				})
				return
			}
			usr = res
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

			res, err := svc.GetUser(claims.UserID)
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

		setErr := svc.SetToken(usr.ID.String(), authToken)
		if setErr != nil {
			log.Error(setErr.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": setErr,
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

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
