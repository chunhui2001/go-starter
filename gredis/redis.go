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

type Redis struct {
	Enable      bool          `mapstructure:"REDIS_Enable"`
	Host        string        `mapstructure:"REDIS_Host"`
	Passwd      string        `mapstructure:"REDIS_Password"`
	MaxIdle     int           `mapstructure:"REDIS_MaxIdle"`
	MaxActive   int           `mapstructure:"REDIS_MaxActive"`
	IdleTimeout time.Duration `mapstructure:"REDIS_IdleTimeout"`
	Db          int           `mapstructure:"REDIS_DataBase"`
}

var (
	redisClient *redis.Client
	ctx         context.Context
	conf        *Redis
	log         *logrus.Logger
)

func Init(redisConf *Redis, log *logrus.Logger) {

	conf = redisConf

	if !redisConf.Enable {
		return
	}

	ctx = context.Background()

	// Connect to Redis
	if redisConf.Passwd != "" {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     redisConf.Host,
			DB:       redisConf.Db,
			Password: redisConf.Passwd,
		})
	} else {
		redisClient = redis.NewClient(&redis.Options{
			Addr: redisConf.Host,
			DB:   redisConf.Db,
		})
	}

	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		log.Error(fmt.Sprintf("Redis client connect failed: Host=%s, errorMessage=%s", redisConf.Host, utils.ErrorToString(err)))
		return
	}

	log.Info(fmt.Sprintf("Redis client connected successfully: Host=%s, DB=%d", redisConf.Host, redisConf.Db))

}

func Client() *redis.Client {
	if !conf.Enable {
		panic(errors.New("Redis-Not-Enabled"))
	}
	return redisClient
}

func Set(key string, value string, expirs time.Duration) {

	err := redisClient.Set(ctx, key, value, expirs).Err()

	if err != nil {
		panic(err)
	}

}
