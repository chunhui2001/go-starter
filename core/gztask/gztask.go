package gztask

import (
	"fmt"
	"time"

	"github.com/chunhui2001/go-starter/core/gzok"
	"github.com/chunhui2001/go-starter/core/utils"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type SimpleGTask struct {
	Enable   bool   `mapstructure:"SIMPLE_GTASK_ENABLE"`
	ID       string `mapstructure:"SIMPLE_GTASK_ID"`
	Name     string `mapstructure:"SIMPLE_GTASK_NAME"`
	Expr     string `mapstructure:"SIMPLE_GTASK_EXPR"`
	PrintLog bool   `mapstructure:"GTASK_PRINT_LOG"`
}

type GZTaskItem struct {
	Id   string
	Name string
	Expr string
}

func (s *GZTaskItem) String() string {
	return fmt.Sprintf(`expr='%s', memo=%s`, s.Expr, s.Name)
}

var (
	nyc, _      = time.LoadLocation("Asia/Shanghai")
	logger      *logrus.Entry
	AppName     string
	currentNode string = utils.OutboundIP().String()
	prefix             = "/__gztask_"
	c                  = cron.New(cron.WithSeconds(), cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)))
	printLog           = true
)

func Init(log *logrus.Entry, appName string, conf *SimpleGTask) {

	logger = log
	AppName = appName
	printLog = conf.PrintLog
	c.Start()

	if conf.Enable {
		AddTask(&GZTaskItem{Id: conf.ID, Expr: conf.Expr, Name: conf.Name}, func() {
			for i := 0; i < 3; i++ {
				time.Sleep(1 * time.Second)
				logger.Infof("GZTask定时任务正在执行每秒1次,耗时3秒: num=%d", i+1)
			}
		})
	}

}

// ### Golang Cron V3 Timed Tasks
// https://www.sobyte.net/post/2021-06/golang-cron-v3-timed-tasks/

// ### 使用在线工具来看自己写的 cron 对不对
// https://en.wikipedia.org/wiki/Cron
// https://crontab.guru/
// "* * * * * *" -- 每秒1次
// "0/5 * * * * *" -- 每5秒
// "0 30 * * * *" -- 每半小时1次
// "15 * * * * *" -- 每15秒1次
// "@hourly" -- Every hour
// "@every 1h30m" -- Every hour thirty
// "@daily" -- Every day

// ## The following is taken from Wikipedia
// # ┌───────────── minute (0 - 59)
// # │ ┌───────────── hour (0 - 23)
// # │ │ ┌───────────── day of the month (1 - 31)
// # │ │ │ ┌───────────── month (1 - 12)
// # │ │ │ │ ┌───────────── day of the week (0 - 6) (Sunday to Saturday;
// # │ │ │ │ │                                   7 is also Sunday on some systems)
// # │ │ │ │ │
// # │ │ │ │ │
// # * * * * * <command to execute>

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

		EntryID, err := c.AddFunc(task.Expr, func() {
			if printLog {
				logger.Infof(`GZTask-Running: lockPath=%s, lockNode=%s, %s`, theLockKey, currentNode, task.String())
			}
			ftask()
		})

		if err != nil {
			logger.Errorf(`GZTask-Add-a-Task-Error: TaskName=%s, TaskId=%s, ErrorMessage=%v`,
				task.Name, AppName+"_"+task.Id, err)
		} else {
			logger.Infof(`GZTask-Add-a-Task: EntryID=%d, lockPath=%s, lockNode=%s, %s`,
				EntryID, theLockKey, currentNode, task.String())
		}

	}()

}
