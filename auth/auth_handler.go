package auth

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"slack-clone-api/domain/user"
	"slack-clone-api/logger"

	"github.com/gin-gonic/gin"
)

func JWTConfigHandler(svc AuthRepository) gin.HandlerFunc {
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
		if reqLogin.GrantType == AuthCode {
			rs, err := validateUser(reqLogin, svc, c)
			if err != nil {
				log.Error(err.Error())
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "username or password is incorrect",
				})
				return
			}
			usr = rs
			n := generateRandomNumber()
			c.JSON(http.StatusOK, gin.H{
				"auth_code": n,
			})
			return

		} else if reqLogin.GrantType == VerifyCode {
			claims, err := clearAuthToken(svc, reqLogin.AuthCode, c)
			if err != nil {
				log.Error(err.Error())
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": err.Error(),
				})
				return
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

		if svc.SetAuthToken(usr.ID.String(), authToken, c); err != nil {
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

func validateUser(reqLogin Login, svc AuthRepository, ctx context.Context) (user.User, error) {
	usr, err := svc.GetUserByEmail(reqLogin.Email, ctx)
	fmt.Println("user", usr)
	if err != nil {
		return usr, err
	}

	return usr, nil
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

func generateRandomNumber() int64 {
	max := big.NewInt(999999)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		fmt.Println(err)
	}

	return n.Int64()
}
