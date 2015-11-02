package main

import (
	"bufio"
	"fmt"
	//	"io/ioutil"
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

	//httpServeMux := http.NewServeMux()
	//httpServeMux.HandleFunc("/recv", recv)
	//httpServeMux.HandleFunc("/push", push)
	//server := &http.Server{
	//	Addr:    ":1234",
	//	Handler: httpServeMux,
	//}
	//server.SetKeepAlivesEnabled(true)
	//server.ListenAndServe()
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
		println("close:", userId)
	}
}

func recvMsg(rwr *bufio.ReadWriter, userId int64) {
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
	//var bodyBytes []byte
	var err error
	params := r.URL.Query()
	suserId := params.Get("userId")
	smsg := params.Get("msg")
	if iuserId, err = strconv.ParseInt(suserId, 10, 16); err != nil {
		fmt.Printf("parse userId is wrong")
		return
	}
	fmt.Printf("%d\t%s\n", iuserId, smsg)
	//bodyBytes, _ := ioutil.ReadAll(r.Body)
	//fmt.Printf("%d\n", len(string(bodyBytes)))
	if _, ok := us.chs[iuserId]; !ok {
		var ta = make(chan Talk, 1)
		var ch Channel
		ch.ch = ta
		us.chs[iuserId] = ch
	}
	talk := Talk{iuserId, smsg}
	us.chs[iuserId].ch <- talk
	//hj, _ := w.(http.Hijacker)
	//conn, _, _ := hj.Hijack()
	//conn.Close()
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
	us = new(Users)
	us.conns = make(map[int64]net.Conn, 200000)
	us.chs = make(map[int64]Channel, 200000)
	initHttp()
}
