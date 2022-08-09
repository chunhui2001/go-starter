package main

import (
	"github.com/chunhui2001/go-starter/starter"

	"github.com/chunhui2001/go-starter/config"
	_ "github.com/chunhui2001/go-starter/cron"
	_ "github.com/chunhui2001/go-starter/gredis"
)

func main() {

	r := starter.Setup()
	r.Run(config.AppSetting.AppPort)

}
