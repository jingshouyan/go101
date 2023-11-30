package main

import (
	_ "go101/logger"
	_ "go101/model"
	"go101/serve"
)

func main() {
	serve.Serve()
}
