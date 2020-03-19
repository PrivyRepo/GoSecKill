package main

import (
	"github.com/gomodule/redigo/redis"
	"io/ioutil"
	"log"
	"net/http"
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
			con, err := redis.Dial("tcp", ":6379")
			if err != nil {
				return nil, err
			}
			return con, nil
		},
	}
}

func check(err error) {
	if err != nil {
		log.Print(err)
	} else {
		log.Print("success")
	}
}

func Print(reply interface{}, e error) {
	log.Println(reply)
	check(e)
}

func main1() {
	conn, e := redis.Dial("tcp", ":6379")
	check(e)
	log.Println("链接成功")
	Print(conn.Do("set", 12, "test"))
	Print(conn.Do("set", []int{1, 2}, "test"))
	Print(conn.Do("get", []int{1, 2}))

}

func main2() {
	conn, e := redis.Dial("tcp", ":6379")
	check(e)
	bytes, i := ioutil.ReadFile("D://README.md")
	check(i)
	reply, e := conn.Do("set", "test", bytes)
	Print(reply, e)
	http.HandleFunc("/testredis", func(rw http.ResponseWriter, r *http.Request) {
		log.Print(r)
		// 从池里获取连接
		rc := redisClient.Get()
		// 用完后将连接放回连接池
		defer rc.Close()
		// 错误判断
		s, e := redis.Bytes(conn.Do("get", "test"))
		check(e)
		rw.Write(s)
	})
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("hello world"))
	})
	e = http.ListenAndServe(":8888", nil)
	check(e)
}
