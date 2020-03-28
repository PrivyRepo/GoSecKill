package redis

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

var redisClient *redis.Pool

func init() {
	// 建立连接池
	redisClient = &redis.Pool{
		MaxIdle:     3,
		MaxActive:   0,
		IdleTimeout: 240 * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			con, err := redis.Dial("tcp", "192.168.19.210:6379")
			if err != nil {
				return nil, err
			}
			return con, nil
		},
	}
}

func NewRedisConn() redis.Conn {
	conn := redisClient.Get()
	return conn
}
