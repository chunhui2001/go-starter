package main

import (
	"github.com/chunhui2001/go-starter/starter"

	"github.com/chunhui2001/go-starter/config"
)

func main() {
	APP_PORT := config.GetEnv("APP_PORT", ":8080")
	starter.Setup().Run(APP_PORT)
}
