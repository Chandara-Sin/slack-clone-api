package auth

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"net/http"
	"slack-clone-api/logger"
	"strconv"

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

		n := generateAuthCode()
		token, _ := svc.InsertAuthToken(strconv.FormatInt(n, 10), c)

		c.JSON(http.StatusOK, gin.H{
			"auth_code":    n,
			"access_token": token,
		})

		// usr := user.User{}
		// if reqLogin.GrantType == AuthCode {
		// usr = rs
		// } else if reqLogin.GrantType == VerifyCode {
		// 	claims, err := clearAuthToken(svc, reqLogin.AuthCode, c)
		// 	if err != nil {
		// 		log.Error(err.Error())
		// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		// 			"error": err.Error(),
		// 		})
		// 		return
		// 	}

		// 	res, err := svc.GetUser(claims.UserID, c)
		// 	if err != nil {
		// 		log.Error(err.Error())
		// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		// 			"error": err.Error(),
		// 		})
		// 		return
		// 	}
		// 	usr = res
		// }

		// authToken, err := GenerateJWTPair(usr)
		// if err != nil {
		// 	log.Error(err.Error())
		// 	c.JSON(http.StatusInternalServerError, gin.H{
		// 		"error": err,
		// 	})
		// 	return
		// }

		// if svc.SetAuthToken(usr.ID.String(), authToken, c); err != nil {
		// 	log.Error(err.Error())
		// 	c.JSON(http.StatusInternalServerError, gin.H{
		// 		"error": err,
		// 	})
		// 	return
		// }

		// c.JSON(http.StatusOK, gin.H{
		// 	"access_token":  authToken.AccessToken,
		// 	"refresh_token": authToken.RefreshToken,
		// 	"token_type":    "Bearer",
		// })
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

		code, err := svc.GetToken(authCode.Token, c)
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

		_, err := clearAuthToken(svc, signOut.Token, c)
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

func clearAuthToken(svc AuthRepository, rfToken string, ctx context.Context) (*JwtCustomClaims, error) {
	token, err := ValidateToken(rfToken)
	if err != nil {
		return nil, err
	}

	claims := GetTokenClaims(token)
	if _, err := svc.GetToken(claims.Subject, ctx); err != nil {
		return nil, errors.New("unauthorized")
	} else {
		svc.ClearToken(claims.UserID, ctx)
		svc.ClearToken(claims.Subject, ctx)
	}

	return claims, nil
}

func generateAuthCode() int64 {
	max := big.NewInt(999999)
	n, _ := rand.Int(rand.Reader, max)
	return n.Int64()
}
