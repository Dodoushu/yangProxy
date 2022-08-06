package main

import (
	"code.byted.org/demo/goPractice/services"
	"log"
	"os"
)
import kingpin "gopkg.in/alecthomas/kingpin.v2"

var (
	app     *kingpin.Application
	service *services.ServiceItem
)

func initConfig() (err error) {

	args := services.Args{}
	tcpArgs := services.TCPArgs{}
	httpArgs := services.HTTPArgs{}

	app = kingpin.New("proxy", "starting with proxy")
	app.Author("yangchengshu").Version("0.0.1")

	args.Parent = app.Flag("parent", "parent address, such as: \"23.32.32.19:28008\"").Default("").Short('P').String()
	args.Local = app.Flag("local", "local ip:port to listen").Short('p').Default(":33080").String()

	//########http#########
	http := app.Command("http", "proxy on http mode")
	httpArgs.LocalType = http.Flag("local-type", "parent protocol type <tls|tcp>").Default("tcp").Short('t').Enum("tls", "tcp")
	httpArgs.ParentType = http.Flag("parent-type", "parent protocol type <tls|tcp>").Short('T').Enum("tls", "tcp")
	httpArgs.Always = http.Flag("always", "always use parent proxy").Default("false").Bool()
	httpArgs.Timeout = http.Flag("timeout", "tcp timeout milliseconds when connect to real server or parent proxy").Default("2000").Int()
	httpArgs.HTTPTimeout = http.Flag("http-timeout", "check domain if blocked , http request timeout milliseconds when connect to host").Default("3000").Int()
	httpArgs.Interval = http.Flag("interval", "check domain if blocked every interval seconds").Default("10").Int()
	httpArgs.Blocked = http.Flag("blocked", "blocked domain file , one domain each line").Default("blocked").Short('b').String()
	httpArgs.Direct = http.Flag("direct", "direct domain file , one domain each line").Default("direct").Short('d').String()
	httpArgs.AuthFile = http.Flag("auth-file", "http basic auth file,\"username:password\" each line in file").Short('F').String()
	httpArgs.Auth = http.Flag("auth", "http basic auth username and password, mutiple user repeat -a ,such as: -a user1:pass1 -a user2:pass2").Short('a').Strings()
	httpArgs.PoolSize = http.Flag("pool-size", "conn pool size , which connect to parent proxy, zero: means turn off pool").Short('L').Default("20").Int()
	httpArgs.CheckParentInterval = http.Flag("check-parent-interval", "check if proxy is okay every interval seconds,zero: means no check").Short('I').Default("3").Int()

	//########tcp#########
	tcp := app.Command("tcp", "proxy on tcp mode")
	tcpArgs.Timeout = tcp.Flag("timeout", "tcp timeout milliseconds when connect to real server or parent proxy").Short('t').Default("2000").Int()
	tcpArgs.ParentType = tcp.Flag("parent-type", "parent protocol type <tls|tcp|udp>").Short('T').Enum("tls", "tcp", "udp")
	tcpArgs.IsTLS = tcp.Flag("tls", "proxy on tls mode").Default("false").Bool()
	tcpArgs.PoolSize = tcp.Flag("pool-size", "conn pool size , which connect to parent proxy, zero: means turn off pool").Short('L').Default("20").Int()
	tcpArgs.CheckParentInterval = tcp.Flag("check-parent-interval", "check if proxy is okay every interval seconds,zero: means no check").Short('I').Default("3").Int()

	kingpin.MustParse(app.Parse(os.Args[1:]))

	httpArgs.Args = args
	tcpArgs.Args = args

	serviceName := kingpin.MustParse(app.Parse(os.Args[1:]))
	services.Regist("http", services.NewHTTP(), httpArgs)
	service, err = services.Run(serviceName)
	if err != nil {
		log.Fatalf("run service [%s] fail, ERR:%s", service, err)
	}
	return
}
