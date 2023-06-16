package gzok

import (
	"time"

	"strings"

	"github.com/go-zookeeper/zk"
	"github.com/sirupsen/logrus"
)

type GZokConf struct {
	Enabled  bool     `mapstructure:"ZOOKEEPER_ENABLED"`
	Debugger bool     `mapstructure:"ZOOKEEPER_DEBUG"`
	Hosts    []string `mapstructure:"ZOOKEEPER_HOSTS"`
}

var (
	logger         *logrus.Entry
	Conn           *zk.Conn
	connectTimeout = time.Duration(5) * time.Second
	lockTimeout    = time.Duration(5) * time.Second
	debugger       = true
)

func Init(zokConf *GZokConf, log *logrus.Entry) {

	logger = log
	debugger = zokConf.Debugger
	conn, _, err := zk.Connect(zokConf.Hosts, connectTimeout)

	if err != nil {
		logger.Errorf(`Zookeeper-Connect-Failed: ConnectTimeout=%s, Servers=%s, ErrorMessage=%s`, connectTimeout, strings.Join(zokConf.Hosts, ","), err.Error())
		return
	}

	Conn = conn

	logger.Infof(`Zookeeper-Connect-Successful: ConnectTimeout=%s, Servers=%s, SessionId=%d`, connectTimeout, strings.Join(zokConf.Hosts, ","), Conn.SessionID())

}

func TryLock(path string, data string) {

	thePath := path

	if !strings.HasPrefix(path, "/") {
		thePath = "/" + path
	}

	for !tryLock(thePath, data) {
		if debugger {
			logger.Debugf(`Zookeeper-TryLock-Failed-Retry: LockTimeout=%s, LockPath=%s, Data=%s`, lockTimeout, thePath, data)
		}
	}

}

// flags 有4种取值：
// 0:永久，除非手动删除
// zk.FlagEphemeral = 1:短暂，session断开则该节点也被删除
// zk.FlagSequence  = 2:会自动在节点后面添加序号
// 3:Ephemeral和Sequence，即，短暂且自动添加序号
func tryLock(path string, data string) bool {

	_, err := Conn.Create(path, []byte(data), zk.FlagEphemeral, zk.WorldACL(zk.PermAll))

	if err == nil {
		logger.Infof(`Zookeeper-TryLock-Successful: LockPath=%s, Data=%s`, path, data)
		return true
	}

	_, _, watch, err := Conn.ExistsW(path)

	if err != nil {
		logger.Errorf(`Zookeeper-Watch-Error: LockPath=%s, ErrorMessage=%s`, path, err.Error())
	}

	select {
	case event := <-watch:
		if event.Type == zk.EventNodeDeleted {
			_, err := Conn.Create(path, []byte(data), zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
			if err == nil {
				logger.Warnf(`Zookeeper-TryLock-Successful: LockPath=%s, Data=%s`, path, data)
				return true
			}
		}
	// time out
	case <-time.After(5 * time.Second):
		return false
	}

	return false

}

func Read(conn *zk.Conn, path string) (string, error) {
	data, _, err := conn.Get(path)
	return string(data), err
}

// 删改与增不同在于其函数中的version参数,其中version是用于 CAS支持, 可以通过此种方式保证原子性
func Modify(conn *zk.Conn, path string, data string) error {
	new_data := []byte(data)
	_, sate, _ := conn.Get(path)
	_, err := conn.Set(path, new_data, sate.Version)
	return err
}

func Delete(conn *zk.Conn, path string) error {
	_, sate, _ := conn.Get(path)
	err := conn.Delete(path, sate.Version)
	return err
}
