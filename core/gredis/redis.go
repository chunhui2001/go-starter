package gredis

import (
	"context"
	"fmt"
	"regexp"
	_ "strconv"
	"strings"
	"time"

	"github.com/chunhui2001/go-starter/core/utils"
	"github.com/go-errors/errors"
	"github.com/go-redis/redis/v8"
	"github.com/gobuffalo/events"
	"github.com/sirupsen/logrus"
)

type REDIS_Mode int64

const (
	Disabled REDIS_Mode = iota
	Standalone
	Sentinel
	Cluster
)

func (s REDIS_Mode) String() string {
	switch s {
	case Disabled:
		return "Disabled"
	case Standalone:
		return "Standalone"
	case Sentinel:
		return "Sentinel"
	case Cluster:
		return "Cluster"
	}
	return "unknown"
}

func (r *GRedis) Disabled() bool {
	return r.Mode == Disabled
}

type MessageHandler func(channel string, payload string)

type GRedis struct {
	Mode           REDIS_Mode    `mapstructure:"REDIS_MODE"` // 0: disable, 1:single, 2:sentinel, 3:cluster
	Host           string        `mapstructure:"REDIS_HOST"`
	Addrs          string        `mapstructure:"REDIS_ADDRS"`
	MasterName     string        `mapstructure:"REDIS_MASTER_NAME"`
	Passwd         string        `mapstructure:"REDIS_PASSWORD"`
	Db             int           `mapstructure:"REDIS_DATABASE"`
	MaxIdle        int           `mapstructure:"REDIS_MAX_IDLE"`
	MaxActive      int           `mapstructure:"REDIS_MAX_ACTIVE"`
	IdleTimeout    time.Duration `mapstructure:"REDIS_IDLE_TIMEOUT"`
	RouteByLatency bool          `mapstructure:"REDIS_ROUTE_BY_LATENCY"`
	RouteRandomly  bool          `mapstructure:"REDIS_ROUTE_RANDOMLY"`
	SubChannels    string        `mapstructure:"REDIS_SUB_CHANNELS"`
	PrintMessage   bool          `mapstructure:"REDIS_MESSAGE_PRINT"`
}

func (r *GRedis) ServerAddrs() string {
	if r.Mode == Standalone {
		return r.Host
	}
	if r.Mode == Sentinel {
		return r.Addrs
	}
	if r.Mode == Cluster {
		return r.Addrs
	}
	return ""
}

var (
	universalClient redis.UniversalClient
	redisClient     *redis.Client
	redisCluster    *redis.ClusterClient
	ctx             context.Context
	conf            *GRedis
	logger          *logrus.Entry
	connected       bool
)

// opt, err := redis.ParseURL("redis://<user>:<pass>@localhost:6379/<db>")
func Init(redisConf *GRedis, log *logrus.Entry) {

	conf = redisConf
	logger = log

	// events.NamedListen("my-listener", func(e events.Event) {
	// 	logger.Infof("### e1 -> %s", e)
	// })

	if conf.Mode == Disabled {
		return
	}

	ctx = context.Background()

	events.Emit(events.Event{
		Kind:    "gredis:Init:start",
		Message: "hi!",
		Payload: events.Payload{"context": ctx},
	})

	// Connect to Redis
	if conf.Mode == Standalone {
		if conf.Passwd != "" {
			redisClient = redis.NewClient(&redis.Options{
				Addr:     conf.Host,
				DB:       conf.Db,
				Password: conf.Passwd,
			})
		} else {
			redisClient = redis.NewClient(&redis.Options{
				Addr: conf.Host,
				DB:   conf.Db,
			})
		}
		universalClient = redisClient
	} else if conf.Mode == Sentinel || conf.Mode == Cluster {

		var addrs []string = strings.Split(conf.Addrs, ",")

		if conf.Mode == Sentinel {

			if conf.RouteByLatency || conf.RouteRandomly {
				if conf.Passwd != "" {
					redisCluster = redis.NewFailoverClusterClient(&redis.FailoverOptions{
						MasterName:     conf.MasterName,
						SentinelAddrs:  addrs,
						RouteByLatency: conf.RouteByLatency,
						RouteRandomly:  conf.RouteRandomly,
						Password:       conf.Passwd,
					})
				} else {
					redisCluster = redis.NewFailoverClusterClient(&redis.FailoverOptions{
						MasterName:     conf.MasterName,
						SentinelAddrs:  addrs,
						RouteByLatency: conf.RouteByLatency,
						RouteRandomly:  conf.RouteRandomly,
					})
				}
				universalClient = redisCluster
			} else {
				if conf.Passwd != "" {
					redisClient = redis.NewFailoverClient(&redis.FailoverOptions{
						MasterName:    conf.MasterName,
						SentinelAddrs: addrs,
						Password:      conf.Passwd,
					})
				} else {
					redisClient = redis.NewFailoverClient(&redis.FailoverOptions{
						MasterName:    conf.MasterName,
						SentinelAddrs: addrs,
					})
				}
				universalClient = redisClient
			}

		} else if conf.Mode == Cluster {
			if conf.RouteByLatency || conf.RouteRandomly {
				redisCluster = redis.NewClusterClient(&redis.ClusterOptions{
					Addrs:          addrs,
					RouteByLatency: conf.RouteByLatency,
					RouteRandomly:  conf.RouteRandomly,
				})
				universalClient = redisCluster
			} else {
				redisCluster = redis.NewClusterClient(&redis.ClusterOptions{
					Addrs: addrs,
				})
				universalClient = redisCluster
			}
		}
	}

	serverVersion := Ping()

	if connected {
		InitProducer(serverVersion)
	}

}

func Ping() string {

	var serverInfo string = "N/a"

	if conf.Mode == Sentinel {
		serverInfo = fmt.Sprintf("Mode=%s, MasterName=%s, ServerAddrs=%s", conf.Mode.String(), conf.MasterName, conf.ServerAddrs())
	} else if conf.Mode == Standalone {
		serverInfo = fmt.Sprintf("Mode=%s, ServerAddrs=%s, DB=%d", conf.Mode.String(), conf.ServerAddrs(), conf.Db)
	} else if conf.Mode == Cluster {
		serverInfo = fmt.Sprintf("Mode=%s, ServerAddrs=%s", conf.Mode.String(), conf.ServerAddrs())
	}

	info, _ := universalClient.Info(context.Background(), "server").Result()

	var serverVersion string = "N/a"
	var redisVersionRE = regexp.MustCompile(`redis_version:(.+)`)

	match := redisVersionRE.FindAllStringSubmatch(info, -1)

	if len(match) < 1 {
		// could not extract redis version
		// ..
	} else {
		version := strings.TrimSpace(match[0][1])
		serverVersion = version
	}

	if _, err := universalClient.Ping(ctx).Result(); err != nil {
		logger.Error(fmt.Sprintf("Redis-Client-Connect-Failed: ServerVersion=%s, %s, errorMessage=%s", serverVersion, serverInfo, utils.ErrorToString(err)))
		connected = false
		return ""
	}

	if conf.SubChannels != "" {
		logger.Info(fmt.Sprintf("Redis-Client-Connected-Successfully: ServerVersion=%s, %s", serverVersion, serverInfo))
	} else {
		logger.Info(fmt.Sprintf("Redis-Client-Connected-Successfully: ServerVersion=%s, %s", serverVersion, serverInfo) + ", SubChannels=" + conf.SubChannels)
	}

	connected = true

	return serverVersion

}

func Client() redis.UniversalClient {
	if conf == nil || conf.Mode == Disabled {
		panic(errors.New("Redis-Not-Enabled"))
	}
	return universalClient
}

func Expire(key string, expiration int) {
	if err := Client().Expire(ctx, key, time.Duration(expiration)*time.Second).Err(); err != nil {
		panic(err)
	}
}

func Del(key ...string) {
	if err := Client().Del(ctx, key...).Err(); err != nil {
		panic(err)
	}
}

func Ttl(key string) (int64, error) {

	val, err := Client().TTL(ctx, key).Result()

	switch {
	case err == redis.Nil:
		return 0, nil
	case err != nil:
		logger.Errorf(`Redis-Get-Key-Error: Key=%s, ErrorMessage=%s`, key, err.Error())
		return 0, err
	}

	return val.Nanoseconds(), nil

}

func Exists(key string) (bool, error) {

	val, err := Client().Exists(ctx, key).Result()

	switch {
	case err == redis.Nil:
		return false, nil
	case err != nil:
		logger.Errorf(`Redis-Get-Key-Error: Key=%s, ErrorMessage=%s`, key, err.Error())
		return false, err
	}

	return val != 0, nil

}

func Get(key string) string {

	val, err := Client().Get(ctx, key).Result()

	switch {
	case err == redis.Nil:
		return ""
	case err != nil:
		logger.Errorf(`Redis-Get-Key-Error: Key=%s, ErrorMessage=%s`, key, err.Error())
		panic(err)
	case val == "":
		return ""
	}

	return val

}

// expir 0 代表无过期时间, 过期时间单位是秒
func Set(key string, value string, expir int) {
	if err := Client().Set(ctx, key, value, time.Duration(expir)*time.Second).Err(); err != nil {
		panic(err)
	}
}

func SetNX(key string, value string, expir int) bool {
	if result, err := Client().SetNX(ctx, key, value, time.Duration(expir)*time.Second).Result(); result {
		return result
	} else {
		if err != nil {
			logger.Errorf(`Redis-SetNX-Error: Key=%s, ErrorMessage=%s`, key, err.Error())
			panic(err)
		}
	}
	return false
}

// 将给定 key 的值设为 value ，并返回 key 的旧值(old value)。
// 当 key 存在但不是字符串类型时，返回一个错误。
// 当 key 没有旧值时，也即是，key 不存在时，返回 null 的同时将当前key设置为新值
func GetSet(key string, value string) string {

	val, err := Client().GetSet(ctx, key, value).Result()

	switch {
	case err == redis.Nil:
		return ""
	case err != nil:
		logger.Errorf(`Redis-GetSet-Key-Error: Key=%s, ErrorMessage=%s`, key, err.Error())
		panic(err)
	case val == "":
		return ""
	}

	return val

}

// 查询列表元素索引,没找到返回-1
// The command returns the index of matching elements inside a Redis list.
// maxLen: 最多找几个
func LindexOf(key string, value string, maxLen int64) int64 {
	val, err := Client().LPos(ctx, key, value, redis.LPosArgs{Rank: 0, MaxLen: maxLen}).Result()
	switch {
	case err == redis.Nil:
		return -1
	case err != nil:
		logger.Errorf(`Redis-LindexOf-Error: Key=%s, ErrorMessage=%s`, key, err.Error())
		panic(err)
	}
	return val
}

// 列表操作
func Lpush(key string, values ...interface{}) {
	if err := Client().LPush(ctx, key, values...).Err(); err != nil {
		panic(err)
	}
}

// 列表操作
func Rpush(key string, values ...interface{}) {
	if err := Client().RPush(ctx, key, values...).Err(); err != nil {
		panic(err)
	}
}

// 读取列表元素: end=-1, 读取所有
func Lrange(key string, start int64, end int64) []string {

	val, err := Client().LRange(ctx, key, start, end).Result()

	switch {
	case err == redis.Nil:
		return []string{}
	case err != nil:
		logger.Errorf(`Redis-Lrange-Error: Key=%s, ErrorMessage=%s`, key, err.Error())
		return nil
	}

	return val
}

// 删除指定范围的列表元素
// start=100, end=-1, 将第100个之前的全部删除, 即保留100个之后的元素
func Ltrim(key string, start int64, end int64) {
	if err := Client().LTrim(ctx, key, start, end).Err(); err != nil {
		panic(err)
	}
}

func Lpop(key string) string {

	val, err := Client().LPop(ctx, key).Result()

	switch {
	case err == redis.Nil:
		return ""
	case err != nil:
		logger.Errorf(`Redis-Lpop-Error: Key=%s, ErrorMessage=%s`, key, err.Error())
		return ""
	}

	return val

}

func Rpop(key string) string {

	val, err := Client().RPop(ctx, key).Result()

	switch {
	case err == redis.Nil:
		return ""
	case err != nil:
		logger.Errorf(`Redis-Rpop-Error: Key=%s, ErrorMessage=%s`, key, err.Error())
		return ""
	}

	return val

}

func Llen(key string) int64 {

	val, err := Client().LLen(ctx, key).Result()

	switch {
	case err == redis.Nil:
		return 0
	case err != nil:
		logger.Errorf(`Redis-Llen-Error: Key=%s, ErrorMessage=%s`, key, err.Error())
		return 0
	}

	return val

}

func Hset(key string, values ...interface{}) {
	if err := Client().HSet(ctx, key, values...).Err(); err != nil {
		panic(err)
	}
}

func Hsetnx(key string, field string, value interface{}) {
	if err := Client().HSetNX(ctx, key, field, value).Err(); err != nil {
		panic(err)
	}
}

func Hget(key string, field string) string {

	val, err := Client().HGet(ctx, key, field).Result()

	switch {
	case err == redis.Nil:
		return ""
	case err != nil:
		logger.Errorf(`Redis-Hget-Error: Key=%s, ErrorMessage=%s`, key, err.Error())
		return ""
	}

	return val
}

func Hgetall(key string) map[string]string {

	val, err := Client().HGetAll(ctx, key).Result()

	switch {
	case err == redis.Nil:
		return nil
	case err != nil:
		logger.Errorf(`Redis-Hgetall-Error: Key=%s, ErrorMessage=%s`, key, err.Error())
		return nil
	}

	return val
}

func Hvals(key string) []string {

	val, err := Client().HVals(ctx, key).Result()

	switch {
	case err == redis.Nil:
		return nil
	case err != nil:
		logger.Errorf(`Redis-Hvals-Error: Key=%s, ErrorMessage=%s`, key, err.Error())
		return nil
	}

	return val

}

func Zincr(key string) (int64, error) {

	result, err := Client().Incr(ctx, key).Result()

	if err != nil {
		logger.Errorf(`Redis-Zincr-Error: Key=%s, ErrorMessage=%s`, key, err.Error())
		return 0, err
	}

	return result, nil

}

func Pub(channel string, payload string) {

	if conf.Mode == Disabled {
		panic(errors.New("Redis-Not-Enabled"))
	}

	err := Client().Publish(ctx, channel, payload).Err()

	if err != nil {
		logger.Error(fmt.Sprintf("Redis-Publish-Error: channel=%s, errorMessage=%s", channel, utils.ErrorToString(err)))
	}

}

func Sub(channel string, handler MessageHandler) {

	if channel == "" {
		return
	}

	if conf != nil && conf.Mode == Disabled {
		panic(errors.New("Redis-Not-Enabled"))
	}

	if !connected {
		logger.Info("Redis-Not-Connected: connected=" + utils.ToString(connected))
		return
	}

	var pubSub *redis.PubSub

	if redisClient != nil {
		pubSub = redisClient.Subscribe(ctx, channel)
	} else if redisCluster != nil {
		pubSub = redisCluster.Subscribe(ctx, channel)
	} else {
		panic(errors.New("Redis-Client-Not-Initializable"))
	}

	// defer pubSub.Close()

	logger.Info("Redis-Subscribe-A-Channel: channel=" + channel)

	go LoopMessage(pubSub, channel, handler)

}

func LoopMessage(pubSub *redis.PubSub, channel string, handler MessageHandler) {

	for {

		msg, err := pubSub.ReceiveMessage(ctx)

		if err != nil {
			logger.Error(fmt.Sprintf("Redis-ReceiveMessage-Error: channel=%s, errorMessage=%s", channel, utils.ErrorToString(err)))
		} else {
			if handler == nil {
				logger.Info("Redis-ReceivedMessage: channel=" + msg.Channel + ", payload=" + msg.Payload)
			} else {
				go handler(channel, msg.Payload)
			}
		}

	}

}
