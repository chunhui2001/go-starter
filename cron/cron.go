package cron

import (
	"github.com/robfig/cron"
)

var c = cron.New()

func init() {

	c.Start()

}

// "* * * * * *" 每秒执行一次
func Add(exp string, f func()) {
	_ = c.AddFunc(exp, f)
}
