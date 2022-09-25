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
	nodeid int64
)

func Init(log *logrus.Entry, node int64) {
	logger = log
	nodeid = node
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
func AddTask(appName string, taskId string, memo string, expr string, task func(node string, taskId string)) {
	logger.Infof(`GRTask-Add-a-Task: taskId=%s, expr='%s', memo=%s`, taskId, expr, memo)
	_ = c.AddFunc(expr, func() {
		lockKey := "__GRTASK_" + appName + "_" + taskId
		Lock(lockKey, taskId, memo, expr, task)
	})
}

func Lock(lockKey string, taskId string, memo string, expr string, task func(node string, taskId string)) {

	currentNode := utils.ToString(nodeid)

	if ok, e := gredis.Exists(lockKey); ok {
		if ttl, err := gredis.Ttl(lockKey); err == nil {
			if ttl <= 0 {
				gredis.Del(lockKey)
			}
			// in progress, 正在执行
			return
		} else {
			logger.Errorf(`GRTask-Ttl-Error: LockKey=%s, expr='%s', ErrorMessage=%s`, lockKey, expr, utils.ErrorToString(err))
			return
		}
	} else {
		if e != nil {
			logger.Errorf(`GRTask-Exists-Error: LockKey=%s, expr='%s', ErrorMessage=%s`, lockKey, expr, utils.ErrorToString(e))
		}
	}

	if gredis.SetNX(lockKey, currentNode, 5) {

		logger.Infof(`GRTask-Started: currentNode=%s, OutboundIP=%s, LockKey=%s`, currentNode, utils.OutboundIP().String(), lockKey)

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
		task(currentNode, lockKey)

		time.Sleep(70 * time.Millisecond) // 暂停70毫秒, 避免定时任务执行的太快, 同时拿到锁

		gredis.Del(lockKey)

		return

	} else {
		logger.Infof(`GRTask-Locked-Discard: LockKey=%s, expr='%s', currentNode=%s, lockVal=%s`, lockKey, expr, currentNode, currentNode)
	}

}
