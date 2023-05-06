package auth

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"slack-clone-api/logger"

	"github.com/gin-gonic/gin"
)

func SignUpHandler(svc AuthRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.Unwrap(c)

		reqSignUp := SignUp{}
		if err := c.ShouldBindJSON(&reqSignUp); err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		_, err := svc.InsertUserByEmail(reqSignUp.Email, c)
		if err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		n := generateAuthCode(6)
		token, _ := svc.SetAuthToken(n, c)

		c.JSON(http.StatusOK, gin.H{
			"auth_code":    n,
			"access_token": token,
		})
	}
}

func AuthCodeHandler(svc AuthRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.Unwrap(c)

		authCode := AuthCode{}
		if err := c.ShouldBindJSON(&authCode); err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		code, err := svc.GetAuthCode(authCode.Token, c)
		if err != nil {
			log.Error(err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		if code != authCode.Code {
			log.Error("wrong auth code")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	}
}

func SignOutHandler(svc AuthRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.Unwrap(c)

		signOut := SignOut{}
		if err := c.ShouldBindJSON(&signOut); err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		err := svc.ClearAuthCode(signOut.Token, c)
		if err != nil {
			log.Error(err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	}
}

func generateAuthCode(maxDigits uint32) string {
	bi, _ := rand.Int(
		rand.Reader,
		big.NewInt(int64(math.Pow(10, float64(maxDigits)))),
	)
	return fmt.Sprintf("%0*d", maxDigits, bi)
}
