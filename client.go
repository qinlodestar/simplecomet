package main

import (
	"fmt"
	"net/http"
	"time"
)

func get(i int) {
	//url := fmt.Sprintf("http://127.0.0.1:8070/sub?ver=1&op=7&seq=1&cb=callback&t=%d", i)
	url := fmt.Sprintf("http://127.0.0.1:1234/recv?userId=%d", i)
	//fmt.Printf("%s\n", url)
	if _, err := http.Get(url); err != nil {
		fmt.Printf("%d\n%v\n", i, err)
	}
}
func main() {
	for i := 0; i <= 50000; i++ {
		if i%50 == 0 {
			time.Sleep(time.Second * 1)
		}
		go get(i)
	}
	time.Sleep(time.Second * 10000)
}
