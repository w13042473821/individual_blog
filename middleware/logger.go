package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		log.Printf("请求: %s %s %d %v", c.Request.Method, c.Request.URL.Path, c.Writer.Status(), duration)
	}
}
