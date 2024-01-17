package utils

import (
	"errors"
	"fmt"
	"log"
	"runtime/debug"
	"sync"
)

// GoPanic 并发调用服务，每个handler执行完后响应
func GoPanic(handlers ...func() error) (err error) {

	var wg sync.WaitGroup
	// 假设我们要调用handlers这么多个服务
	for _, f := range handlers {

		wg.Add(1)
		// 每个函数启动一个协程
		go func(handler func() error) {
			defer func() {
				// 每个协程内部使用recover捕获可能在调用逻辑中发生的panic
				if e := recover(); e != nil {
					// 某个服务调用协程报错，可以在这里打印一些错误日志
					log.Println(e)
					debug.PrintStack()
				}
				wg.Done()
			}()

			// 取第一个报错的handler调用逻辑，并最终向外返回
			e := handler()
			if err == nil && e != nil {
				err = e
			}
		}(f)
	}

	wg.Wait()
	return
}

// GoPanicToError 并发调用服务，只要一个 handler panic 程序 shutdown
func GoPanicToError(handlers func() error) error {
	stopChan := make(chan string)
	// 假设我们要调用handlers这么多个服务
	// for _, f := range handlers {
	// 每个函数启动一个协程
	go func(handler func() error) {
		defer func() {
			// 每个协程内部使用recover捕获可能在调用逻辑中发生的panic
			e := recover()
			if e != nil {
				log.Println(e)
				//打印错误堆栈信息
				debug.PrintStack()
				stopChan <- fmt.Sprintf("%s", e)
			}

		}()
		// 取第一个报错的handler调用逻辑，并最终向外返回
		err := handler()
		if err != nil {
			stopChan <- err.Error()
		}
	}(handlers)
	// }

	t := <-stopChan
	return errors.New(t)
}
