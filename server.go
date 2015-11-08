package main

import (
	"bufio"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"net"
	"net/http"
	"strconv"
)

func initHttp() {
	http.HandleFunc("/recv", recv)
	http.HandleFunc("/push", push)
	http.HandleFunc("/", test)
	fmt.Printf("starting server...\n")
	http.ListenAndServe(":1234", nil)
}

// create pool
func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "127.0.0.1:6379", redis.DialPassword("moodecn2015"))
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
	ok, err := redis.String(c.Do("SET", key, "127.0.0.1:1234"))
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

func test(w http.ResponseWriter, r *http.Request) {
	cn, _ := w.(http.CloseNotifier)
	go closeHandler(1, cn)
}

func closeHandler(userId int64, cn http.CloseNotifier) {
	<-cn.CloseNotify()
}

func recv(w http.ResponseWriter, r *http.Request) {
	var iuserId int64
	var err error
	params := r.URL.Query()
	suserId := params.Get("userId")
	if iuserId, err = strconv.ParseInt(suserId, 10, 64); err != nil {
		fmt.Printf("parse userId is wrong\n")
		return
	}
	fmt.Printf("recv=%d\n", iuserId)

	hj, _ := w.(http.Hijacker)
	conn, rwr, _ := hj.Hijack()
	us.conns[iuserId] = conn

	go recvMsg(rwr, iuserId)
	go read(rwr, iuserId)
}

/* read client data, if client closed,then close server conn
 *
 */
func read(rwr *bufio.ReadWriter, userId int64) {
	var p byte = '\n'
	var err error
	if _, err = rwr.ReadString(p); err != nil {
		conn := us.conns[userId]
		if closeErr := conn.Close(); closeErr != nil {
			println("close error:", userId)
		}
		delete(us.conns, userId)
		del(userId)
		println("close:", userId)
	}
}

func recvMsg(rwr *bufio.ReadWriter, userId int64) {
	var (
		ok       bool
		conn     net.Conn
		closeErr error
	)

	//close old connection and channel
	if _, ok = us.chs[userId]; ok {
		conn = us.conns[userId]
		if closeErr = conn.Close(); closeErr != nil {
			println("close old connection error:", userId)
		}
	}

	set(userId)

	res, _ := get(userId)
	fmt.Printf("%v\n", res)

	var ta = make(chan Talk, 1)
	var ch Channel
	ch.ch = ta
	us.chs[userId] = ch
	for {
		talk := <-us.chs[userId].ch
		fmt.Printf("recv:%v\n", talk)
		str := fmt.Sprintf("<script>alert(\"%d:%s\")</script>", talk.userId, talk.msg)
		rwr.WriteString(str)
		rwr.Flush()
	}
}

func push(w http.ResponseWriter, r *http.Request) {
	var iuserId int64
	var err error
	params := r.URL.Query()
	suserId := params.Get("userId")
	smsg := params.Get("msg")
	if iuserId, err = strconv.ParseInt(suserId, 10, 16); err != nil {
		fmt.Printf("parse userId is wrong")
		return
	}
	fmt.Printf("%d\t%s\n", iuserId, smsg)
	if _, ok := us.chs[iuserId]; !ok {
		var ta = make(chan Talk, 1)
		var ch Channel
		ch.ch = ta
		us.chs[iuserId] = ch
	}
	talk := Talk{iuserId, smsg}
	us.chs[iuserId].ch <- talk
}

type Users struct {
	conns map[int64]net.Conn
	chs   map[int64]Channel
}

type Channel struct {
	ch chan Talk
}

type Talk struct {
	userId int64
	msg    string
}

var us *Users

// create redis connection pool
var pool = newPool()

func main() {
	us = new(Users)
	us.conns = make(map[int64]net.Conn, 200000)
	us.chs = make(map[int64]Channel, 200000)
	initHttp()
}
