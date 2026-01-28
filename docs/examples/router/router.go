package router

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/docs/examples/controller"
)

func Main() {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 注册
	new(controller.DemoBindingController).Register(r)

	if err := r.Run(":8080"); err != nil {
		fmt.Printf("error:%v\n", err)
		os.Exit(1)
	}
}
