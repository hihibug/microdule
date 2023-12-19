package web

import "net/http"

type Web interface {
	Client() *Gin
	Init(...Option) error
	Options() Options
	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
	Stop() error
	Run() error
}

// Handler 回调
type Handler func(Event) error

// Message 消息
type Message struct {
	Header map[string]string
	Body   []byte
}

// Event 事件
type Event interface {
	Topic() string
	Message()
}

var DefauktWeb = NewGin()
