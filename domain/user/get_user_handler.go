package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type getUserFunc func(string) (User, error)

func (fn getUserFunc) GetUser(ID string) (User, error) {
	return fn(ID)
}

func GetUserHanlder(svc getUserFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ID := c.Param("id")

		usr, err := svc.GetUser(ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, usr)
	}
}
