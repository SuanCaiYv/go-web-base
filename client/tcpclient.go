package client

import (
	"fmt"
	"log"
	"net"
)

func Work() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8190")
	if err != nil {
		log.Fatal(err)
	}
	tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatal(err)
	}
	var input string
	response := make([]byte, 1024)
	fmt.Println("input for 5 loops")
	for i := 0; i < 5; i++ {
		_, _ = fmt.Scanf("%s", &input)
		_, _ = tcpConn.Write([]byte(input))
		readLen, _ := tcpConn.Read(response)
		fmt.Println(string(response[:readLen]))
	}
}
