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
//
// 2022/01/24更 今天看了看一个比较有名的WebSocket包的实现，大致了解了WebSocket的原理：
// 首先我们需要知道Golang的http包处理请求的基本流程是读取Socket，包装成Request和Response，然后作为传参调用我们自己编写的处理函数。
// 当然，调用处理函数这一步已经在新的Go程中进行了；Golang的Response对象，实现了Hijack()方法，这个方法会获取当前HTTP请求背后的Socket。
// 正常情况下，调用业务处理函数之后，http包会帮我们调用关闭和刷新等方法来实现收尾操作，但是如果调用hijack()之后；
// 缓冲区的刷新和释放，连接的关闭等事件就需要我们自己来完成；http包将不再替我们收尾。
// 此时我们得到了Socket之后，读取请求头，判断是否是Upgrade请求，如果是，则接管Socket，然后我们自行写入对应的升级响应，顺带一提，HTTP判断响应结束是基于分隔符的。
// 所以我们简单地追加分隔符(即手动构建响应)即可实现响应完成的通知，但是此时socket还在；后续的操作就完全基于这个Socket啦！实现全双工。
// 以上就是大致原理，其实很简单，完全可以自行实现WebSocket。
// 读取请求 => 是否是升级请求 => 接管Socket => 写出升级响应 => 不关闭Socket连接 => 全双工读写 => 执行收尾。
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
