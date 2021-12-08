package server

import (
	"github.com/gorilla/websocket"
	"net/http"
)

// WsServe WebSocket协议弥补了HTTP的不足——无法实现全双工通信；它的基本工作原理如下：
// 它和HTTP一样，都是需要建立连接的，所以为了建立
// 安全连接，它需要先握手，而握手的方式考虑到了兼容性和省事，直接用了HTTP Get请求来完成。
// 这个Get请求的请求头包含了一些字段，比如URI(即服务端收到对这个URI的访问，就知道是客户端想升级到WebSocket了)；
// 再比如加密字段，服务端收到之后解密，并追加一个魔法值(固定的值)，然后加密发送给客户端，这是为了保证服务端也支持WS，而不会
// 把此次Get请求当成一个HTTP给吃了。
// 客户端收到Get的响应后，判断如果响应符合预期，就直接走Socket建立全双工通信了。
// 其实握手的过程也可以看成打开Socket的过程，只是写法换成了另一种，而这种方法恰好通过HTTP来实现比较简单。
// 除了一些必要字段，还可以的设置数据传输协议，比如json，然后客户端和服务端通信都会以json格式为基础进行数据传输
func WsServe() {
}

var upgrader = websocket.Upgrader{
	HandshakeTimeout:  0,
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	WriteBufferPool:   nil,
	Subprotocols:      nil,
	Error:             nil,
	CheckOrigin:       nil,
	EnableCompression: false,
}

func handle(writer http.ResponseWriter, request *http.Request) {
	connection, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		return
	}
	connection.ReadMessage()
}
