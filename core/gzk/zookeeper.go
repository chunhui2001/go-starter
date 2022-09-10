package gzk

import (
	"strings"
	"time"

	"context"

	DLocker "github.com/nladuo/go-zk-lock"
	"github.com/olebedev/emitter"
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
		NewLock("mylock2", 12, func(locked bool) {
			if locked {
				logger.Infof(`Zookeeper-Get-Simple-Lock-Succeed: %s, executed`, "FocusLock")
			}
		}).AcquireLock()
	}

}

type Lock struct {
	Event   *emitter.Emitter
	Path    string
	TimeOut int
}

func NewLock(path string, timeOut int, f func(bool)) *Lock {

	e := &emitter.Emitter{}

	lock := &Lock{
		Event:   e,
		Path:    path,
		TimeOut: timeOut,
	}

	go func() {
		for event := range e.On("change") {
			status := event.Int(0) // cast the first argument to int
			if status == 2 {
				f(true)
			} else {
				time.Sleep(5 * time.Second)
				lock.AcquireLock() // 超时重拿
			}
		}
	}()

	return lock
}

func (l *Lock) AcquireLock() {

	currLockPath := chroot + "_" + l.Path

	go func() {

		mylock := make(chan DLocker.Dlocker, 1)
		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*time.Duration(l.TimeOut))
		defer cancelFunc()

		go newLocker(ctx, mylock, currLockPath)

		go func() {
			select {
			case <-time.After(time.Duration(l.TimeOut) * time.Second):
				// time out
				logger.Errorf(`Zookeeper-Get-Locked-Error-TimeOut: Path=%s, Error=%s`, currLockPath, ctx.Err())
				l.Event.Emit("change", 1)
			case <-ctx.Done():
				logger.Warnf(`Zookeeper-Get-Locked-Context-Done: Path=%s, Error=%s`, currLockPath, ctx.Err())
				l.Event.Emit("change", 3)
			case <-mylock:
				logger.Infof(`Zookeeper-Get-Locked-Succeed: Path=%s, Status=%s`, currLockPath, ctx.Err())
				// l.F(true) // like mutex.Lock()
				l.Event.Emit("change", 2)
			}
		}()

	}()

}

func newLocker(ctx context.Context, mylock chan DLocker.Dlocker, path string) {
	locker := DLocker.NewLocker(path, time.Duration(1)*time.Hour) // 锁100年
	locker.Lock()
}

func FocusLock(lockPath string, timeOut int, f func(bool)) {

	currLockPath := chroot + lockPath

	go func() {

		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*time.Duration(timeOut))

		defer cancelFunc()

		go func() {
			select {
			case <-time.After(time.Duration(timeOut) * time.Second):
				// time out
				if deadline, ok := ctx.Deadline(); ok {
					logger.Warnf(`Zookeeper-Get-Locked-Failed: Path=%s, Deadline=%s`, currLockPath, time.Since(deadline))
					f(false)
				} else {
					logger.Errorf(`Zookeeper-Get-Locked-Error-TimeOut: Path=%s, Deadline=%s, Error=%s`, currLockPath, time.Since(deadline), ctx.Err())
					f(false)
				}
			case <-ctx.Done():
				if deadline, ok := ctx.Deadline(); ok {
					logger.Infof(`Zookeeper-Get-Locked-Succeed: Path=%s, SpentTime=%s, Status=%s`, currLockPath, time.Until(deadline), ctx.Err())
					f(true) // like mutex.Lock()
				} else {
					logger.Warnf(`Zookeeper-Get-Locked-Error-Done: Path=%s, Deadline=%s, Error=%s`, currLockPath, time.Since(deadline), ctx.Err())
					f(false)
				}
			}
		}()

		locker := DLocker.NewLocker(currLockPath, time.Duration(999999)*time.Hour) // 锁100年
		locker.Lock()

		ctx.Done()

	}()

}
