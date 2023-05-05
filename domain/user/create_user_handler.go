package user

import (
	"context"
	"net/http"
	"slack-clone-api/logger"

	"github.com/gin-gonic/gin"
)

type createUserFunc func(User, context.Context) error

func (fn createUserFunc) CreateUser(usr User, ctx context.Context) error {
	return fn(usr, ctx)
}

func CreateUserHanlder(svc createUserFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.Unwrap(c)

		usr := User{}
		if err := c.ShouldBindJSON(&usr); err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		err := svc.CreateUser(usr, c)
		if err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"status": "ok",
		})
	}
}
