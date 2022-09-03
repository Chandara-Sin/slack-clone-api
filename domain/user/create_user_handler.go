package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type createUserFunc func(User) error

func (fn createUserFunc) CreateUser(usr User) error {
	return fn(usr)
}

func CreateUserHanlder(svc createUserFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		usr := User{}

		if err := c.ShouldBindJSON(&usr); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		err := svc.CreateUser(usr)
		if err != nil {
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
