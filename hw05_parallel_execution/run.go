package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

var errorCount int
var mu sync.Mutex

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}
	var (
		taskChan = make(chan Task, n)
		// errorChan = make(chan error, m)
		wg        sync.WaitGroup
		globalErr error
	)
	// go func() {
	// 	var errorCount int

	// 	for {
	// 		select {
	// 		case _, ok := <-errorChan:
	// 			if ok {
	// 				errorCount++
	// 			} else {
	// 				return
	// 			}
	// 		default:
	// 			if errorCount >= m {
	// 				globalErr = ErrErrorsLimitExceeded
	// 				return
	// 			}
	// 		}
	// 	}
	// }()
	go func() {
		defer close(taskChan)
		for _, task := range tasks {
			taskChan <- task
		}
	}()

	for i := 0; i < n; i++ {
		wg.Add(1)
		go handler(taskChan, &globalErr, m, &wg)
	}

	wg.Wait()
	return globalErr
}

func handler(taskChan chan Task, globalErr *error, maxErrorCount int, wg *sync.WaitGroup) {
	defer wg.Done()
	defer fmt.Println("DONE")
	for task := range taskChan {
		if *globalErr != nil {
			return
		}
		if err := task(); err != nil {
			mu.Lock()
			errorCount++
			mu.Unlock()
		}
		if errorCount == maxErrorCount {
			*globalErr = ErrErrorsLimitExceeded
			return
		}
	}
}
