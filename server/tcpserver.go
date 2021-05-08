package server

import (
	"fmt"
	"log"
	"net"
	"time"
)

func work(tcpConn *net.TCPConn) {
	defer func(tcpConn *net.TCPConn) {
		err := tcpConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(tcpConn)
	// 设置30秒未读到就关闭本次请求，这里是所有读取的总时长！！！
	_ = tcpConn.SetReadDeadline(time.Now().Add(30 * time.Second))
	request := make([]byte, 1024)
	// 可能客户端在建立了连接后，会请求好几次。
	// 所以此时不能草率的结束。
	for {
		readLen, err := tcpConn.Read(request)
		fmt.Println(readLen)
		if err != nil {
			return
		}
		// 说明连接断了
		if readLen == 0 {
			fmt.Println("connection closed")
			return
		}
		msg := string(request[:readLen])
		if msg == "echo" {
			_, _ = tcpConn.Write([]byte("hello client!"))
		} else if msg == "time" {
			_, _ = tcpConn.Write([]byte(time.Now().String()))
		} else {
			_, _ = tcpConn.Write([]byte("echo: " + msg))
		}
		request = make([]byte, 1024)
	}
}

func work0(tcpConn *net.TCPConn) {
	defer func(tcpConn *net.TCPConn) {
		err := tcpConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(tcpConn)
	request := make([]byte, 1024)
	readLen, _ := tcpConn.Read(request)
	msg := string(request[:readLen])
	fmt.Println(msg)
	msg = "echo: " + msg
	_, _ = tcpConn.Write([]byte(msg))
}

func Serve() {
	address := "127.0.0.1:8190"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		log.Fatal(err)
	}
	// listener对应ServerSocket
	tcpListener, err := net.ListenTCP("tcp4", tcpAddr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		// 每次连接建立返回一个Connection，Connection对应Socket
		tcpConn, err := tcpListener.AcceptTCP()
		fmt.Println("connection established...")
		if err != nil {
			log.Fatal(err)
		}
		go work(tcpConn)
	}
}
