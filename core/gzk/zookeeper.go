package gzk

import (
	"strings"
	"time"

	"context"

	DLocker "github.com/nladuo/go-zk-lock"
	"github.com/sirupsen/logrus"
)

var (
	logger  *logrus.Entry
	servers string               // the zookeeper hosts
	chroot  string = "/__locker" //the application znode path
)

type GZk struct {
	Enable     bool   `mapstructure:"ZOOKEEPER_ENABLE"`
	Servers    string `mapstructure:"ZOOKEEPER_SERVERS"` // the zookeeper hosts
	ChRoot     string `mapstructure:"ZOOKEEPER_CHROOT"`  // the application znode path
	TimeOut    int    `mapstructure:"ZOOKEEPER_TIMEOUT"` // the zk connection timeout // 20 * time.Second
	SimpleLock bool   `mapstructure:"ZOOKEEPER_SIMPLE_LOCK"`
}

func Init(gzk *GZk, log *logrus.Entry) {

	logger = log

	err := DLocker.EstablishZkConn(strings.Split(gzk.Servers, ","), time.Duration(gzk.TimeOut)*time.Second)

	if err != nil {
		logger.Errorf(`Zookeeper-Connect-Failed: Servers=%s, TimeOut=%d/s, ErrorMessage=%s`, gzk.Servers, gzk.TimeOut, err.Error())
		return
	}

	logger.Infof(`Zookeeper-Connected-Successful: TimeOut=%d/s, ChRoot=%s, Servers=%s`, gzk.TimeOut, gzk.ChRoot, gzk.Servers)

	chroot = gzk.ChRoot
	servers = gzk.Servers

	if gzk.SimpleLock {
		FocusLock("lock2", func() {
			logger.Infof(`Zookeeper-Get-Simple-Lock-Succeed: %s, executed`, "FocusLock")
		})
	}

}

func FocusLock(lockPath string, f func()) {

	go func() {

		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*15)
		defer cancelFunc()

		go func() {
			select {
			case <-time.After(2 * time.Second):
				// time out
				logger.Infof(`Zookeeper-Get-Locked-Failed: Path=%s, TimeOut=%s`, chroot+"/"+lockPath, 2*time.Second)
			case <-ctx.Done():
				if deadline, ok := ctx.Deadline(); ok {
					logger.Infof(`Zookeeper-Get-Locked-Succeed: Path=%s, Deadline=%s`, chroot+"/"+lockPath, time.Until(deadline))
				} else {
					logger.Infof(`Zookeeper-Get-Locked-Error: Path=%s, Deadline=%s, Error=%s`, chroot+"/"+lockPath, time.Since(deadline), ctx.Err())
				}
			}
		}()

		locker := DLocker.NewLocker(chroot+"/"+lockPath, time.Duration(999999)*time.Hour) // 锁100年
		locker.Lock()                                                                     // like mutex.Lock()

		ctx.Done()

		f()

	}()

}
