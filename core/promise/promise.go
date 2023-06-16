package promise

import (
	"sync"
	"time"

	"github.com/chunhui2001/go-starter/core/config"
)

var (
	logger = config.Log
)

func WaitGroup(timeOut int, f ...func()) bool {

	var wg sync.WaitGroup
	// var mu sync.Mutex

	wg.Add(len(f))

	for _, fn := range f {

		go func(fn func()) {

			defer func() {
				if r := recover(); r != nil {
					logger.Errorf("WaitGroup-Func-Error: Error=%v", r)
				}
				// mu.Unlock()
				wg.Done()
			}()

			// 使用互斥锁保护对调用者的并发访问
			// mu.Lock()
			fn()

		}(fn)
	}

	done := make(chan struct{})
	timeout := time.After(time.Duration(timeOut) * time.Second)

	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// 所有 Goroutine 执行完毕，没有超时
		return true
	case <-timeout:
		// 超时，打印超时消息
		logger.Warnf(`WaitGroup-Timeout: timeOut=%d`, timeOut)
		return false
	}

}
