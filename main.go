package main

import (
	"log"
	"slack-clone-api/config"
	"slack-clone-api/domain/user"
	"slack-clone-api/store"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig()
	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) {
		c.Status(200)
	})

	db := store.CreateDB()
	if err := db.AutoMigrate(&user.User{}); err != nil {
		log.Println("auto migrate db: ", err)
	}

	u := r.Group("/api")
	u.POST("/users", user.CreateUserHanlder(user.Create(db)))

	r.Run()
}
