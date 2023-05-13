package auth

import (
	"crypto/rand"
	b64 "encoding/base64"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"slack-clone-api/domain/mail"
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
		usr, err := svc.InsertUserByEmail(reqSignUp.Email, c)
		if err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		authCode := generateAuthCode(6)
		err = mail.MailHandler(usr, authCode)
		if err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		token, err := svc.SetAuthToken(usr.Email, authCode, c)
		if err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token": token,
			"token_type":   "Bearer",
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

		if code, err := svc.GetAuthCode(authCode.Token, c); err != nil {
			log.Error(err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		} else if code != authCode.Code {
			log.Error("Invalid Auth Code")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		if err := svc.ClearAuthCode(authCode.Token, c); err != nil {
			log.Error(err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		usr, err := svc.GetUserByEmail(decodeBase64(authCode.Token), c)
		if err != nil {
			log.Error(err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		token, _ := svc.GenerateToken(usr)
		c.JSON(http.StatusOK, gin.H{
			"session_token": token,
			"token_type":    "ID",
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

func decodeBase64(value string) string {
	decoded, _ := b64.StdEncoding.DecodeString(value)
	return string(decoded[:])
}
