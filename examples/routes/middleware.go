package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func middlewareRegistry(engine *gin.Engine) {
	engine.Use(gin.Recovery())
	engine.Use(LogMiddleware())
}

func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("api count")
		c.Next()
	}
}
