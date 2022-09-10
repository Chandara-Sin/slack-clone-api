package login

import (
	"net/http"
	"slack-clone-api/auth"
	"slack-clone-api/domain/user"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type getUserByEmailFunc func(eml string) (user.User, error)

func (fn getUserByEmailFunc) GetUserByEmail(eml string) (user.User, error) {
	return fn(eml)
}

func LoginHanlder(svc getUserByEmailFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqLogin := Login{}

		if err := c.ShouldBindJSON(&reqLogin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		usr, err := svc.GetUserByEmail(reqLogin.Email)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "username or password is incorrect",
			})
			return
		}

		match := checkPasswordHash(reqLogin.Password, usr.HashedPassword)
		if !match {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "username or password is incorrect",
			})
			return
		}

		authToken, err := auth.GenerateJWTPair(usr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token":  authToken.AccessToken,
			"refresh_token": authToken.RefreshToken,
		})
	}
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
