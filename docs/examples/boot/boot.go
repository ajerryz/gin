package boot

import (
	"flag"

	"github.com/gin-gonic/gin/docs/examples/routes"
)

func Main() {

	flag.Parse()

	routes.BootGin()
}
