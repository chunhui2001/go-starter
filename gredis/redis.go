package gredis

import (
	"context"
	"fmt"
	_ "strconv"
	"strings"
	"time"

	"github.com/chunhui2001/go-starter/utils"
	"github.com/go-errors/errors"
	"github.com/go-redis/redis/v8"
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

	if conf.Mode == Disabled {
		return
	}

	ctx = context.Background()

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
				redisCluster = redis.NewFailoverClusterClient(&redis.FailoverOptions{
					MasterName:     conf.MasterName,
					SentinelAddrs:  addrs,
					RouteByLatency: conf.RouteByLatency,
					RouteRandomly:  conf.RouteRandomly,
				})
				universalClient = redisCluster
			} else {
				redisClient = redis.NewFailoverClient(&redis.FailoverOptions{
					MasterName:    conf.MasterName,
					SentinelAddrs: addrs,
				})
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

	Ping()

}

func Ping() {

	var serverInfo string = "N/a"

	if conf.Mode == Sentinel {
		serverInfo = fmt.Sprintf("Mode=%s, MasterName=%s, ServerAddrs=%s", conf.Mode.String(), conf.MasterName, conf.ServerAddrs())
	} else if conf.Mode == Standalone {
		serverInfo = fmt.Sprintf("Mode=%s, ServerAddrs=%s, DB=%d", conf.Mode.String(), conf.ServerAddrs(), conf.Db)
	} else if conf.Mode == Cluster {
		serverInfo = fmt.Sprintf("Mode=%s, ServerAddrs=%s", conf.Mode.String(), conf.ServerAddrs())
	}

	if _, err := universalClient.Ping(ctx).Result(); err != nil {
		logger.Error(fmt.Sprintf("Redis-Client-Connect-Failed: %s, errorMessage=%s", serverInfo, utils.ErrorToString(err)))
		connected = false
		return
	}

	if conf.SubChannels != "" {
		logger.Info(fmt.Sprintf("Redis-Client-Connected-Successfully: %s", serverInfo))
	} else {
		logger.Info(fmt.Sprintf("Redis-Client-Connected-Successfully: %s", serverInfo) + ", SubChannels=" + conf.SubChannels)
	}

	connected = true

}

func Client() redis.UniversalClient {
	if conf == nil || conf.Mode == Disabled {
		panic(errors.New("Redis-Not-Enabled"))
	}
	return universalClient
}

func Set(key string, value string, expir int) {

	err := Client().Set(ctx, key, value, time.Duration(expir)*time.Second).Err()

	if err != nil {
		panic(err)
	}

}

func Get(key string) []byte {

	data, err := Client().Get(ctx, key).Bytes()

	if err != nil {
		panic(err)
	}

	return data

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
