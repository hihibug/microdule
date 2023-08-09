package utils

import (
	"fmt"
	"sync"
)

// GoPanic 并发调用服务，每个handler都会传入一个调用逻辑函数
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
					fmt.Println("并发执行错误")
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
