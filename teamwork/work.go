package teamwork

import (
	"context"
	"errors"
	"log"
	"runtime/debug"

	"golang.org/x/exp/maps"
)

type Teamwork interface {
	Reginster(name string, habdle func()) TeamworkClose
	WorkNum() int
	Start() error
	Close()
}

type TeamworkClose interface {
	HandleClose(func())
}

type TeamworkCloseStruct struct {
	T *TeamworkStruct
}

type TeamworkStruct struct {
	Ctx          context.Context
	Handles      map[string]func()
	HandlesClose []func()
}

func NewTeamwork() Teamwork {
	return &TeamworkStruct{
		Ctx:          context.Background(),
		Handles:      make(map[string]func()),
		HandlesClose: make([]func(), 0),
	}
}

func (t *TeamworkStruct) Reginster(name string, handle func()) TeamworkClose {
	maps.Copy(t.Handles, map[string]func(){name: handle})
	return &TeamworkCloseStruct{
		T: t,
	}
}

func (t *TeamworkStruct) WorkNum() int {
	return len(t.Handles)
}

func (t *TeamworkCloseStruct) HandleClose(handle func()) {
	t.T.HandlesClose = append(t.T.HandlesClose, handle)
}

func (t *TeamworkStruct) Start() error {
	stopChan := make(chan string)
	// 假设我们要调用handlers这么多个服务
	for k, f := range t.Handles {
		// 每个函数启动一个协程
		go func(handler func(), k string) {
			defer func() {
				// 每个协程内部使用recover捕获可能在调用逻辑中发生的panic
				if e := recover(); e != nil {
					log.Println(e)
					//打印错误堆栈信息
					debug.PrintStack()
				}
				stopChan <- "panic shutdown"
			}()
			// 取第一个报错的handler调用逻辑，并最终向外返回
			log.Println(k, "start: Success")
			handler()
		}(f, k)
	}

	err := <-stopChan
	return errors.New(err)
}

func (t *TeamworkStruct) Close() {
	for _, v := range t.HandlesClose {
		v()
	}
}
