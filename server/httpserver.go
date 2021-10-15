package server

import (
	"fmt"
	"net/http"
	"strings"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

type myHandler struct {
	// 我们这里路径映射匹配只在最开始是写的，所以不需要同步
	handlers map[string]HandlerFunc
}

func NewMyHandler() *myHandler {
	return &myHandler{
		handlers: make(map[string]HandlerFunc),
	}
}

func (h *myHandler) AddHandler(path, method string, handler http.Handler) {
	key := path + "#" + method
	h.handlers[key] = handler.ServeHTTP
}

func (h *myHandler) AddHandlerFunc(path, method string, f HandlerFunc) {
	key := path + "#" + method
	h.handlers[key] = f
}

type notFound struct {
}

func (n *notFound) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
}

var handler404 = notFound{}

func (h *myHandler) getHandlerFunc(path, method string) HandlerFunc {
	key := path + "#" + method
	handler, ok := h.handlers[key]
	if !ok {
		// todo 返回404专有handler
		return handler404.ServeHTTP
	} else {
		return handler
	}
}

func (h *myHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	url := request.RequestURI
	method := request.Method
	uri := strings.Split(url, "?")[0]
	h.getHandlerFunc(uri, method)(writer, request)
}

func ServeHttp() {
	myHandler := NewMyHandler()
	myHandler.AddHandlerFunc("/hello", "GET", func(w http.ResponseWriter, r *http.Request) {
		// 必须解析哈！不然会报错
		_ = r.ParseForm()
		fmt.Println(r.Form.Get("name"))
		_, _ = w.Write([]byte("ok"))
	})
	myHandler.AddHandlerFunc("/hello", "POST", func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		fmt.Println(r.PostForm.Get("name"))
		_, _ = w.Write([]byte("ok"))
	})
	myHandler.AddHandlerFunc("/upload", "POST", func(w http.ResponseWriter, r *http.Request) {
		// 限制大小为8MB
		_ = r.ParseMultipartForm(8 << 20)
		fileHeader := r.MultipartForm.File["my_file"][0]
		fmt.Println(fileHeader.Filename)
		_, _ = w.Write([]byte("ok"))
	})
	_ = http.ListenAndServe(":8190", myHandler)
}
