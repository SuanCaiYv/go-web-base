package server

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

// 最普通的版本
func server(socket *net.TCPConn) {
	defer func(tcpConn *net.TCPConn) {
		err := tcpConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(socket)
	for {
		request := make([]byte, 1024)
		readLen, err := socket.Read(request)
		if err == io.EOF {
			fmt.Println("连接关闭")
			return
		}
		msg := string(request[:readLen])
		fmt.Println(msg)
		msg = "echo: " + msg
		_, _ = socket.Write([]byte(msg))
	}
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

// 经过测试，TCP传输数据时，如果一次性写入过大的数据，会使得TCP进行分片传输
// 其实这是基础知识，只是突然想知道TCP所谓的稳定传输到底是什么意思，以及怎么处理TCP的关闭。
// TCP的三次握手是为了确认"远程可达"及"信道无阻"等传输要求，四次挥手是为了实现数据传输扫尾操作，确保不会有数据遗漏，且双方均"认可"关闭这一动作。
// TCP的拥塞控制，滑动窗口，AK确认机制等是为了确保数据一定被目标接收了，如果没有则会触发重传等一系列后备操作。
// 简而言之，TCP协议可以确保：链路存在+数据可达。而UDP只管丢，不管其他的。
// 如果发送方突然挂掉，接收方是感知不到的，因为TCP的可达更多是于发送方而言的。
// 正常关闭情况下，相当于发送了关闭请求，接收方可以知晓这一事件并主动记录关闭操作，说白了这次关闭是双方都知道的"事实"。
// 异常关闭情况下，比如网络中断，宕机等，接收方是感知不到的，因为TCP不会主动告知接收者发送者没了这一"事实"。
// 此时的写出和读取均会超时，因为得不到响应，所以很多语言的库的读写操作允许设置超时就是这个原因。
// TCP本身在发消息时无法保证消息一定会被接收，它只是可以收到此次发送的响应罢了，然后通过响应"得知"成功与否，而不能完全保证发送一定会被收到。
// 所以有了各种可靠机制，比如重传和窗口，都是为了在第一次失败后重新发送的。
// 其实说到这里，我们就可以知道怎么合理处理关闭请求了，如果是主动关闭，则我们这边也一起关闭就行，如果是被动关闭，则设置超时时间进行判断。
func largeTransform(socket *net.TCPConn) {
	defer func(socket *net.TCPConn) {
		_ = socket.Close()
	}(socket)
	reader := bufio.NewReader(socket)
	for {
		// MTU最大不会大于这个值, 64KB
		buf := make([]byte, 1<<22, 1<<22)
		n, _ := reader.Read(buf)
		fmt.Printf("read: %d", n)
		time.Sleep(time.Second)
	}
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
		socket, err := serverSocket.AcceptTCP()
		fmt.Println("connection established...")
		if err != nil {
			log.Fatal(err)
		}
		// 开辟Goroutine去处理新的连接
		go largeTransform(socket)
	}
}
