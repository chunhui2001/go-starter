package main

import (
	"github.com/chunhui2001/go-starter/starter"

	"github.com/chunhui2001/go-starter/config"
)

var APP_PORT string = config.GetEnv("APP_PORT", ":8080")

func main() {

	r := starter.Setup()
	r.Run(APP_PORT)

}
