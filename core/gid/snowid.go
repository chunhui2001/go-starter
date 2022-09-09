package gid

import (
	"github.com/bwmarrin/snowflake"
	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Entry
	gid    *snowflake.Node
)

func Init(log *logrus.Entry, node int64) {

	logger = log
	id, err := snowflake.NewNode(node)

	if err != nil {
		logger.Errorf(`Snowflake-NewNode-Error: Node=%d, ErrorMessage=%s`, node, err.Error())
		return
	}

	gid = id

	logger.Infof(`Snowflake-NewNode-Initialized: Node=%d, RandomId=%d`, node, gid.Generate().Int64())

}

func Get(node int64) int64 {
	return gid.Generate().Int64()
}
