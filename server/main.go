package main

import (
	log "code.google.com/p/log4go"
	"flag"
	"net"
)

var us *Users

func main() {
	flag.Parse()
	if err := InitConfig(); err != nil {
		panic(err)
	}
	//c, _ := json.Marshal(Conf)
	//fmt.Printf("%v\n", string(c))
	log.LoadConfiguration(Conf.Log)
	InitKafka(Conf.KafkaHost)
	popKafka()

	us = new(Users)
	us.conns = make(map[int64]net.Conn, Conf.ConnNum)
	us.chs = make(map[int64]Channel, Conf.ConnNum)
	initHttp()
}
