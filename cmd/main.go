package main

import (
	"gdragon/internal/router"
)

func main() {
	r := router.SetupRouter()
	r.Run(":8090")
}
