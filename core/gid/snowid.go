package gid

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"math"
	"math/big"
	"strings"

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

func Get() int64 {
	return gid.Generate().Int64()
}

func ID() string {
	return YtID()
}

func YtID() string {

	val, err := rand.Int(rand.Reader, big.NewInt(int64(math.MaxInt64)))

	if err != nil {
		panic(err)
	}

	b := make([]byte, 8)

	binary.LittleEndian.PutUint64(b, uint64(val.Int64()))
	encoded := base64.StdEncoding.EncodeToString([]byte(b))

	var replacer = strings.NewReplacer(
		"+", "-",
		"/", "_",
	)

	return replacer.Replace(encoded[:11])

}
