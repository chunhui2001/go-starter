package gztask

import (
	"fmt"
	"github.com/chunhui2001/go-starter/core/gzok"
	"github.com/chunhui2001/go-starter/core/utils"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
)

type GZTaskItem struct {
	Id   string
	Name string
	Expr string
}

func (s *GZTaskItem) String() string {
	return fmt.Sprintf(`expr='%s', memo=%s`, s.Expr, s.Name)
}

var (
	c           = cron.New()
	logger      *logrus.Entry
	AppName     string
	currentNode string = utils.OutboundIP().String()
	prefix             = "/__gztask_"
)

func Init(log *logrus.Entry, appName string) {
	logger = log
	AppName = appName
	c.Start()
}

// "* * * * * *" -- 每秒1次
// "0/5 * * * * *" -- 每5秒
// "0 30 * * * *" -- 每半小时1次
// "15 * * * * *" -- 每15秒1次
// "@hourly" -- Every hour
// "@every 1h30m" -- Every hour thirty
// "@daily" -- Every day

// ###################################################################################
// -----                  | -----------                                | -------------
// Entry                  | Description                                | Equivalent To
// -----                  | -----------                                | -------------
// @yearly (or @annually) | Run once a year, midnight, Jan. 1st        | 0 0 0 1 1 *
// @monthly               | Run once a month, midnight, first of month | 0 0 0 1 * *
// @weekly                | Run once a week, midnight between Sat/Sun  | 0 0 0 * * 0
// @daily (or @midnight)  | Run once a day, midnight                   | 0 0 0 * * *
// @hourly                | Run once an hour, beginning of hour        | 0 0 * * * *
// ###################################################################################
func AddTask(task *GZTaskItem, ftask func()) {

	theLockKey := prefix + AppName + "_" + task.Id

	go func() {

		gzok.TryLock(theLockKey, currentNode)

		logger.Infof(`GZTask-Add-a-Task: lockPath=%s, lockNode=%s, %s`, theLockKey, currentNode, task.String())

		_ = c.AddFunc(task.Expr, func() {
			runTask(theLockKey, task, ftask) // run task
		})

	}()

}

func runTask(theLockKey string, task *GZTaskItem, ftask func()) {
	logger.Infof(`GZTask-Running: lockPath=%s, lockNode=%s, %s`, theLockKey, currentNode, task.String())
	ftask()
}
