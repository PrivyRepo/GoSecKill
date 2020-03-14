package common

import "github.com/gomodule/redigo/redis"

func NewRedisConn() (redis.Conn, error) {
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	return conn, err
}
