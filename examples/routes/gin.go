package routes

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/examples/config"
)

var defaultGinEngine *gin.Engine

func DefaultEngine() *gin.Engine {
	return defaultGinEngine
}

func BootGin() {
	gin.SetMode(*config.GinMode)

	defaultGinEngine = gin.New()

	defaultGinEngine.Use()

	middlewareRegistry(defaultGinEngine) // 注册中间件
	routeRegistry(defaultGinEngine)      // 注册路由

	port := ":" + strconv.Itoa(*config.Port)
	if err := defaultGinEngine.Run(port); err != nil {
		panic(fmt.Errorf("start gin error: %w", err))
	}
}
