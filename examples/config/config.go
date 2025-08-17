package config

import (
	"flag"

	"github.com/gin-gonic/gin"
)

var Port = flag.Int("port", 8080, "Port to listen on")

var GinMode = flag.String("ginMode", gin.ReleaseMode, "ginMode:debug,release,test")
