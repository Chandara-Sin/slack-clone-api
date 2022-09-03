package main

import (
	"fmt"
	"log"
	"slack-clone-api/domain/user"
	"slack-clone-api/store"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	initConfig()
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

func initConfig() {
	viper.SetConfigName("config-dev") // name of config file (without extension)
	viper.SetConfigType("yaml")       // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")          // optionally look for config in the working directory
	err := viper.ReadInConfig()       // Find and read the config file
	if err != nil {
		fmt.Printf("Fatal error config file: %s \n", err) // Handle errors reading the config file
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}
