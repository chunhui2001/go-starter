package main

import (
	"go-starter/config"
	"go-starter/starter"
)

func main() {
	APP_PORT := config.GetEnv("APP_PORT", ":8080")
	starter.Setup().Run(APP_PORT)
}
