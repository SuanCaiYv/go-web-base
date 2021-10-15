package main

import (
	"../server"
	"math"
	"time"
)
import "../client"

func main() {
	server.ServeHttp()
	go server.Serve()
	time.Sleep(10 * time.Millisecond)
	go client.Work()
	time.Sleep(time.Duration(math.MaxUint32) * time.Second)
}
