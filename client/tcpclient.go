package client

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

func clientDelimiterBased() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8190")
	if err != nil {
		log.Fatal(err)
	}
	socket, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatal(err)
	}
	var input string
	fmt.Println("input for 5 loops")
	for i := 0; i < 5; i++ {
		_, _ = fmt.Scanf("%s", &input)
		// 添加分隔符
		input = input + "\n"
		_, _ = socket.Write([]byte(input))
		response := make([]byte, 1024)
		readLen, _ := socket.Read(response)
		fmt.Println(string(response[:readLen]))
	}
	err = socket.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func clientLengthBased() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8190")
	if err != nil {
		log.Fatal(err)
	}
	socket, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatal(err)
	}
	var input string
	fmt.Println("input for 5 loops")
	for i := 0; i < 5; i++ {
		_, _ = fmt.Scanf("%s", &input)
		data := []byte(input)
		var buffer = bytes.NewBuffer([]byte{})
		// 先写入长度
		_ = binary.Write(buffer, binary.BigEndian, int32(len(data)))
		// 再写入数据
		_ = binary.Write(buffer, binary.BigEndian, data)
		_, _ = socket.Write(buffer.Bytes())
		response := make([]byte, 1024)
		readLen, _ := socket.Read(response)
		fmt.Println(string(response[:readLen]))
	}
	err = socket.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func Work() {
	clientLengthBased()
}
