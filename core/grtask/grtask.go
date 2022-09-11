package grtask

import (
	"time"

	"github.com/chunhui2001/go-starter/core/gredis"
	"github.com/chunhui2001/go-starter/core/utils"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Entry
	c      = cron.New()
)

func Init(log *logrus.Entry) {
	logger = log
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
func AddTask(taskId string, memo string, expr string, task func(taskId string)) {
	logger.Infof(`GRTask-Add-a-Task: taskId=%s, expr='%s', memo=%s`, taskId, expr, memo)
	_ = c.AddFunc(expr, func() {
		lockKey := "__GRTASK_APP_NAME_" + taskId
		Lock(lockKey, taskId, memo, expr, task)
	})
}

func Lock(lockKey string, taskId string, memo string, expr string, task func(taskId string)) {

	currentNode := utils.Hostname() + "/" + utils.OutboundIP().String()

	if ok, _ := gredis.Exists(lockKey); ok {
		if ttl, err := gredis.Ttl(lockKey); err == nil {
			if ttl.String() == "-1ns" {
				gredis.Del(lockKey)
			} else {
				return
			}
		} else {
			return
		}
	}

	if gredis.SetNX(lockKey, currentNode, 5) {

		lockedNode := gredis.Get(lockKey)

		logger.Infof(`GRTask-Run-Task-Started: LockKey=%s, expr='%s', lockedNode=%s`, lockKey, expr, lockedNode)

		// 避免定时任务执行时间过长给当前锁续命，避免重复启动
		go func() {
			for {
				time.Sleep(100 * time.Millisecond)
				if ok, _ := gredis.Exists(lockKey); ok {
					gredis.Set(lockKey, currentNode, 5) // 安保线程, 里边的人没出来外边的人进不去
				} else {
					break
				}
			}
		}()

		// 拿到了
		task(lockKey)
		gredis.Del(lockKey)
		return

	}

}
