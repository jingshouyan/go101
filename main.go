package main

import (
	_ "go101/logger"
	_ "go101/model"
	"go101/router"
)

func main() {
	router.Serve()
}
