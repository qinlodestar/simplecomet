simplecomet
==============

simplecomet 是基于golang开发的一个简易的长连接服务器，目前支持只支持http协议，可以兼容任意浏览器。

---------------------------------------

* install 
```sh
go get -u github.com/qinlodestar/simplecomet
cd $GOPATH/src/github.com/qinlodestar/simplecomet
go build
./simplecomet -c server.conf
```
安装过程中会遇到log4go 的日志不能安装，我是先下载log4go,然后再新建 code.google.com/p/log4go为顺序的目录。
* use 

在浏览器上打开
http://127.0.0.1:1234/recv?userId=123

在另外一个页面打开
http://127.0.0.1:1234/push?userId=123&msg=abc

127.0.0.1是部署服务器的地址，在第一个页面就会弹出一个alert，哈哈，通信结束，之后就可以在浏览器上通过长连接接收消息。
