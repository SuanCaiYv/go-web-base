package main

import (
	"go-web-base/client"
	"go-web-base/server"
	"math"
	"time"
)

func main() {
	// server.ServeHttp()
	go server.Serve()
	time.Sleep(10 * time.Millisecond)
	go client.Work()
	time.Sleep(time.Duration(math.MaxUint32) * time.Second)
}
