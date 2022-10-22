package main

import (
	"github.com/chunhui2001/go-starter/core/config"
	"github.com/chunhui2001/go-starter/core/googleapi"
)

func main() {
	googleapi.Init(config.GoogleAPIConfSettings, config.Log)
}
