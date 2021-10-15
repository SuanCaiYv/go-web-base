package server

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

// 最普通的版本
func server(tcpConn *net.TCPConn) {
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

// 基于分隔符的版本
func serverClientDelimiterBased(socket *net.TCPConn) {
	defer func(socket *net.TCPConn) {
		err := socket.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(socket)
	// 构建一个Reader，此时会源源不断的读取，直到Socket为空
	reader := bufio.NewReader(socket)
	for {
		// 相当于对源源不断的数据流进行分割，直到不可读取
		data, err := reader.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				// 连接关闭
				break
			} else {
				fmt.Println("出现异常" + err.Error())
			}
		}
		// 剔除分隔符
		data = data[:len(data)-1]
		text := string(data)
		fmt.Println("服务端读到了: " + text)
		resp := fmt.Sprintf("Hello, client. I have read: [%s] from you.", text)
		_, _ = socket.Write([]byte(resp))
	}
	fmt.Println("连接关闭")
}

func serverClientLengthBased(socket *net.TCPConn) {
	defer func(socket *net.TCPConn) {
		_ = socket.Close()
	}(socket)
	reader := bufio.NewReader(socket)
	for {
		// 先读取长度字段(但是不移动字节指针)
		lenData, err := reader.Peek(4)
		if err != nil {
			// 客户端关闭了连接
			if err == io.EOF {
				break
			} else {
				fmt.Println("出现异常" + err.Error())
			}
		}
		lenBuf := bytes.NewBuffer(lenData)
		var length int32
		_ = binary.Read(lenBuf, binary.BigEndian, &length)
		// 说明只是读到了请求的一部分，还有剩下的请求需要在下次读取，索性全部留到下次循环处理
		if int32(reader.Buffered())-4 < length {
			continue
		}
		data := make([]byte, length+4)
		_, _ = reader.Read(data)
		text := string(data[4:])
		fmt.Println("服务端读到了: " + text)
		resp := fmt.Sprintf("Hello, client. I have read: [%s] from you.", text)
		_, _ = socket.Write([]byte(resp))
	}
	fmt.Println("连接关闭")
}

func Serve() {
	address := "127.0.0.1:8190"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		log.Fatal(err)
	}
	// listener对应ServerSocket
	serverSocket, err := net.ListenTCP("tcp4", tcpAddr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		// 每次连接建立返回一个Connection，Connection对应Socket
		tcpConn, err := serverSocket.AcceptTCP()
		fmt.Println("connection established...")
		if err != nil {
			log.Fatal(err)
		}
		// 开辟Goroutine去处理新的连接
		go serverClientLengthBased(tcpConn)
	}
}
