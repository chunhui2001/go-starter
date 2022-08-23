package gredis

import (
	"context"
	"fmt"
	_ "strconv"
	"time"

	"github.com/chunhui2001/go-starter/utils"
	"github.com/go-errors/errors"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type GRedis struct {
	Enable         bool          `mapstructure:"REDIS_Enable"`
	Host           string        `mapstructure:"REDIS_Host"`
	Addrs          string        `mapstructure:"REDIS_ADDRS"`
	MasterName     string        `mapstructure:"REDIS_MASTER_NAME"`
	Passwd         string        `mapstructure:"REDIS_Password"`
	Db             int           `mapstructure:"REDIS_DataBase"`
	MaxIdle        int           `mapstructure:"REDIS_MaxIdle"`
	MaxActive      int           `mapstructure:"REDIS_MaxActive"`
	IdleTimeout    time.Duration `mapstructure:"REDIS_IdleTimeout"`
	RouteByLatency bool          `mapstructure:"REDIS_RouteByLatency"`
	RouteRandomly  bool          `mapstructure:"REDIS_RouteRandomly"`
}

type Cmdable func(ctx context.Context, cmd redis.Cmder) error

var (
	redisClient  *redis.Client
	redisCluster *redis.ClusterClient
	ctx          context.Context
	conf         *GRedis
	log          *logrus.Logger
)

// opt, err := redis.ParseURL("redis://<user>:<pass>@localhost:6379/<db>")
func Init(redisConf *GRedis, log *logrus.Logger) {

	conf = redisConf

	if !conf.Enable {
		return
	}

	ctx = context.Background()

	// Connect to Redis
	if conf.Addrs != "" {
		// sentinel or cluster
		if conf.MasterName != "" {
			if conf.RouteByLatency || conf.RouteRandomly {
				redisCluster = redis.NewFailoverClusterClient(&redis.FailoverOptions{
					MasterName:     conf.MasterName,
					SentinelAddrs:  []string{":7000", ":7001", ":7002", ":7003", ":7004", ":7005"},
					RouteByLatency: conf.RouteByLatency,
					RouteRandomly:  conf.RouteRandomly,
				})
			} else {
				redisClient = redis.NewFailoverClient(&redis.FailoverOptions{
					MasterName:    conf.MasterName,
					SentinelAddrs: []string{":7000", ":7001", ":7002", ":7003", ":7004", ":7005"},
				})
			}
		} else {
			if conf.RouteByLatency || conf.RouteRandomly {
				redisCluster = redis.NewClusterClient(&redis.ClusterOptions{
					Addrs:          []string{":7000", ":7001", ":7002", ":7003", ":7004", ":7005"},
					RouteByLatency: conf.RouteByLatency,
					RouteRandomly:  conf.RouteRandomly,
				})
			} else {
				redisCluster = redis.NewClusterClient(&redis.ClusterOptions{
					Addrs: []string{":7000", ":7001", ":7002", ":7003", ":7004", ":7005"},
				})
			}
		}

	} else {
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

	}

	if redisClient != nil {
		if _, err := redisClient.Ping(ctx).Result(); err != nil {
			log.Error(fmt.Sprintf("Redis client connect failed: Host=%s, errorMessage=%s", conf.Host, utils.ErrorToString(err)))
			return
		}
	} else if redisCluster != nil {
		if _, err := redisCluster.Ping(ctx).Result(); err != nil {
			log.Error(fmt.Sprintf("Redis client connect failed: Host=%s, errorMessage=%s", conf.Host, utils.ErrorToString(err)))
			return
		}
	}

	log.Info(fmt.Sprintf("Redis client connected successfully: Host=%s, DB=%d", conf.Host, conf.Db))

}

func Ping() {
	log.Info("Redis client connected successfully")
}

func Client() *redis.Client {
	if !conf.Enable {
		panic(errors.New("Redis-Not-Enabled"))
	}
	return redisClient
}

func Set(key string, value string, expir int) {

	err := redisClient.Set(ctx, key, value, time.Duration(expir)*time.Second).Err()

	if err != nil {
		panic(err)
	}

}
