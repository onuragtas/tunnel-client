package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tunnel-client/models"
	request2 "tunnel-client/request"
)

var defaultLocalPort = "80"
var defaultDestinationIp = "127.0.0.1"

var requestClient = request2.Request{
	BaseUrl: "http://tunnel.resoft.org/api/v1",
}

var domainList = models.Domain{}

type StartedTunnels struct {
	Data []Item
}

var startedTunnels StartedTunnels
var process string
var list string

func updateULimit() {
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		fmt.Println("Error Getting Rlimit ", err)
	}
	fmt.Println(rLimit)
	rLimit.Max = 999999
	rLimit.Cur = 999999
	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		fmt.Println("Error Setting Rlimit ", err)
	}
	err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		fmt.Println("Error Getting Rlimit ", err)
	}
	fmt.Println("Rlimit Final", rLimit)
}

func init() {
	updateULimit()
	flag.StringVar(&process, "process", "", "")
	flag.StringVar(&list, "list", "", "")

	flag.Parse()
	SetupCloseHandler()
	go keepAlive()
}

func keepAlive() {
	for true {
		for key, value := range startedTunnels.Data {
			go func(startedTunnels StartedTunnels, item Item, key int) {
				if item.Domain.KeepAlive != 0 && time.Now().Sub(item.KeepAliveTime).Seconds() > float64(item.Domain.KeepAlive) {
					startedTunnels.Data[key].KeepAliveTime = time.Now()
					resp, _ := http.Get("https://" + item.Domain.Domain)
					if resp != nil && resp.Body != nil {
						resp.Body.Close()
					}
				}
			}(startedTunnels, value, key)
		}
		time.Sleep(time.Millisecond * 200)
	}
}

func SetupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Will Closing")
		Close()
		os.Exit(1)
	}()
}

func Close() {
	fmt.Println("Closed...")
}
