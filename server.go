package main

import (
	"bufio"
	log "code.google.com/p/log4go"
	//	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"strconv"
)

func initHttp() {
	http.HandleFunc("/recv", recv)
	http.HandleFunc("/push", push)
	http.HandleFunc("/", test)
	log.Debug("starting server...")
	http.ListenAndServe(Conf.HttpBind, nil)
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
		log.Debug("parse userId is wrong")
		return
	}
	log.Debug("recv=%d", iuserId)

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
	status := "Ok"
	if _, err = rwr.ReadString(p); err != nil {
		conn := us.conns[userId]
		if closeErr := conn.Close(); closeErr != nil {
			status = "Fail"
		}
		delete(us.conns, userId)
		del(userId)
		log.Debug("close%s:\t%d", status, userId)
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
	log.Debug("%v", res)

	var ta = make(chan Talk, 1)
	var ch Channel
	ch.ch = ta
	us.chs[userId] = ch
	for {
		talk := <-us.chs[userId].ch
		log.Debug("recv:%v", talk)
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
		log.Debug("parse userId is wrong")
		return
	}
	log.Debug("%d\t%s", iuserId, smsg)
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

func main() {
	flag.Parse()
	if err := InitConfig(); err != nil {
		panic(err)
	}
	//c, _ := json.Marshal(Conf)
	//fmt.Printf("%v\n", string(c))
	log.LoadConfiguration(Conf.Log)
	us = new(Users)
	us.conns = make(map[int64]net.Conn, Conf.ConnNum)
	us.chs = make(map[int64]Channel, Conf.ConnNum)
	initHttp()
}
