package user

import (
	"context"
	"net/http"
	"slack-clone-api/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type getUserFunc func(uuid.UUID, context.Context) (User, error)

func (fn getUserFunc) GetUser(ID uuid.UUID, ctx context.Context) (User, error) {
	return fn(ID, ctx)
}

func GetUserHanlder(svc getUserFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.Unwrap(c)

		user, _ := c.Get("user")
		ID := user.(*User).ID

		usr, err := svc.GetUser(ID, c)
		if err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, usr)
	}
}
