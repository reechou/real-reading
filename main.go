package main

import (
	"github.com/reechou/real-reading/config"
	"github.com/reechou/real-reading/controller"
)

func main() {
	controller.NewLogic(config.NewConfig()).Run()
}
