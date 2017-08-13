package main

import (
	"flag"

	"github.com/fuckFE/haishi_server/server"
)

func main() {
	port := flag.String("port", "8000", "tcp port")

	flag.Parse()
	r := server.GetMainEngine()

	r.Run(":" + *port)
}
