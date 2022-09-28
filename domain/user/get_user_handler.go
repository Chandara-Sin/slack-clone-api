package user

import (
	"net/http"
	"slack-clone-api/logger"

	"github.com/gin-gonic/gin"
)

type getUserFunc func(string) (User, error)

func (fn getUserFunc) GetUser(ID string) (User, error) {
	return fn(ID)
}

func GetUserHanlder(svc getUserFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.Unwrap(c)

		user, _ := c.Get("user")
		ID := user.(*User).Id

		usr, err := svc.GetUser(ID.String())
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
