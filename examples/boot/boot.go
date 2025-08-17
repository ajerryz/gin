package boot

import (
	"flag"

	"github.com/gin-gonic/gin/examples/routes"
)

func Main() {

	flag.Parse()

	routes.BootGin()
}
