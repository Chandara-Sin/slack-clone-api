package logger

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const key = "logger"

func Middleware(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.Request.Header.Get("X-Request-ID")
		l := log.With(zap.String("req-id", reqID))
		c.Set(key, l)

		// Read the Body content
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
		}

		// Restore the io.ReadCloser to its original state
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		fmt.Printf("request body: %s\n", bodyBytes)

		c.Next()
	}
}

func Unwrap(c gin.Context) *zap.Logger {
	val, _ := c.Get(key)
	if log, ok := val.(*zap.Logger); ok {
		return log
	}
	return zap.NewExample()
}
