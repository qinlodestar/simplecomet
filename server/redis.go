package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

// create redis connection pool
var pool = newPool()

// create pool
func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", Conf.RedisHost, redis.DialPassword(Conf.RedisPwd))
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

func set(userId int64) bool {
	c := pool.Get()
	defer c.Close()
	key := fmt.Sprintf("comet:%d", userId)
	ok, err := redis.String(c.Do("SETEX", key, 86400, "127.0.0.1:1234"))
	if ok != "OK" || err != nil {
		return false
	}
	return true
}

func get(userId int64) (string, error) {
	c := pool.Get()
	defer c.Close()
	key := fmt.Sprintf("comet:%d", userId)
	res, err := redis.String(c.Do("GET", key))
	return res, err
}

func del(userId int64) bool {
	c := pool.Get()
	defer c.Close()
	key := fmt.Sprintf("comet:%d", userId)
	ok, err := redis.String(c.Do("DELETE", key))
	if ok != "OK" || err != nil {
		return false
	}
	return true
}
