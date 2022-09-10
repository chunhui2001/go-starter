package gzk

import (
	"strings"
	"time"

	"github.com/go-zookeeper/zk"
	"github.com/sirupsen/logrus"
)

var (
	conn *zk.Conn
)

func Init2(gzk *GZk, log *logrus.Entry) {

	logger = log
	c, _, err := zk.Connect(strings.Split(gzk.Servers, ","), time.Duration(10)*time.Second)

	if err != nil {
		logger.Errorf(`Zk-Connect-Failed: Servers=%s, TimeOut=%d/s, ErrorMessage=%s`, gzk.Servers, gzk.TimeOut, err.Error())
		return
	}

	logger.Infof(`Zookeeper-Connected-Successful: TimeOut=%d/s, ChRoot=%s, Servers=%s`, gzk.TimeOut, gzk.ChRoot, gzk.Servers)

	conn = c
	AcquireLock2("/__lock666", func(locked bool) {
		logger.Infof(`Zk-Get-Simple-Lock-Succeed: %s, executed`, "FocusLock")
	})
}

func AcquireLock2(path string, f func(bool)) {
	for {
		if ok, _, _ := conn.Exists(path); ok {
			time.Sleep(5 * time.Second)
			continue
		}
		conn.Create(path, []byte(""), int32(0), zk.WorldACL(zk.PermAll))
		// children, stat, ch, err := conn.Create(path, []byte(""), int32(0), zk.WorldACL(zk.PermAll))
		logger.Infof(`Zk-Get-Locked-Succeed: Path=%s`, path)
		f(true)
	}
}
