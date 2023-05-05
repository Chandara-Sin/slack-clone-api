package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"slack-clone-api/auth"
	"slack-clone-api/config"
	"slack-clone-api/domain/user"
	"slack-clone-api/logger"
	"slack-clone-api/mw"
	"slack-clone-api/store"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	config.InitConfig()
	r := gin.Default()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	zaplog, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer zaplog.Sync()
	r.Use(ginzap.Ginzap(zaplog, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(zaplog, true))
	r.Use(logger.Middleware(zaplog))

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:8000",
		"http://localhost:3000",
	}
	config.AllowHeaders = append(
		config.AllowHeaders,
		"Authorization",
		"X-API-KEY",
	)
	config.AllowCredentials = true
	r.Use(cors.New(config))

	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "Ok v1")
	})

	db := store.CreateDB()
	rdb := store.InitRedisDB(ctx)

	a := r.Group("/api")
	a.Use(mw.ValidatorOnlyAPIKey(viper.GetString("api.key.public")))

	authRepository := auth.AuthRepository{
		DB:  db,
		RDB: rdb,
	}
	a.POST("/oauth/signup", auth.LoginHandler(authRepository))
	a.POST("/users", user.CreateUserHanlder(user.Create(db)))

	u := r.Group("/api")
	u.Use(mw.JWTConfig(viper.GetString("jwt.secret"), mw.GetToken(rdb)))
	u.GET("/users/info", user.GetUserHanlder(user.GetUser(db)))
	u.POST("/oauth/revoke", auth.SignOutHandler(authRepository))

	s := &http.Server{
		Addr:           ":" + viper.GetString("app.port"),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	stop()
	fmt.Println("shutting down gracefully, press Ctrl+C again to force")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.Shutdown(timeoutCtx); err != nil {
		fmt.Println(err)
	}
}
