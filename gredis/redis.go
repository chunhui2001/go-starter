package gredis

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/chunhui2001/go-starter/config"
	"github.com/chunhui2001/go-starter/logger"
	"github.com/gomodule/redigo/redis"
)

var RedisConn *redis.Pool

// Setup Initialize the Redis instance
func init() {

	REDIS_Enable, _ := strconv.ParseBool(config.GetEnv("REDIS_Enable", "false"))

	if REDIS_Enable == false {
		return
	}

	REDIS_HOST := config.GetEnv("REDIS_Host", "127.0.0.1:6379")
	REDIS_Password := config.GetEnv("REDIS_Password", "")
	REDIS_MaxIdle, _ := strconv.Atoi(config.GetEnv("REDIS_MaxIdle", "30"))
	REDIS_MaxActive, _ := strconv.Atoi(config.GetEnv("REDIS_MaxActive", "30"))
	REDIS_IdleTimeout, _ := time.ParseDuration(config.GetEnv("REDIS_IdleTimeout", "200"))

	RedisConn = &redis.Pool{

		MaxIdle:     REDIS_MaxIdle,
		MaxActive:   REDIS_MaxActive,
		IdleTimeout: REDIS_IdleTimeout,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", REDIS_HOST)
			if err != nil {
				logger.Log.Error("connect to redis error: errorMessage=" + fmt.Sprint(err))
				return nil, err
			}
			if REDIS_Password != "" {
				if _, err := c.Do("AUTH", REDIS_Password); err != nil {
					c.Close()
					logger.Log.Error("connect to redis error: errorMessage=" + fmt.Sprint(err))
					return nil, err
				}
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			logger.Log.Error("connect to redis error: errorMessage=" + fmt.Sprint(err))
			return err
		},
	}

	return
}

// Set a key/value
func Set(key string, data interface{}, time int) error {
	conn := RedisConn.Get()
	defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", key, value)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, time)
	if err != nil {
		return err
	}

	return nil
}

// Exists check a key
func Exists(key string) bool {
	conn := RedisConn.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return exists
}

// Get get a key
func Get(key string) ([]byte, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

// Delete delete a kye
func Delete(key string) (bool, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

// LikeDeletes batch delete
func LikeDeletes(key string) error {
	conn := RedisConn.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	for _, key := range keys {
		_, err = Delete(key)
		if err != nil {
			return err
		}
	}

	return nil
}
