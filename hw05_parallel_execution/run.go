package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}
	var (
		taskChan  chan Task  = make(chan Task, n)
		errorChan chan error = make(chan error, m)
		wg        sync.WaitGroup
		globalErr error = nil
	)

	go func() {
		var errorCount int64 = 0
		for {
			if _, ok := <-errorChan; ok {
				errorCount++
			}
			if errorCount >= int64(m) {
				close(errorChan)
				return
			}
		}
	}()

	go func() {
		defer close(taskChan)
		for _, task := range tasks {
			taskChan <- task
		}
	}()

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			handler(taskChan, errorChan, &globalErr)
			wg.Done()
		}()
	}

	wg.Wait()

	return globalErr
}

func handler(taskChan chan Task, errorChan chan error, globalErr *error) {
	if *globalErr != nil {
		return
	}
	for task := range taskChan {
		if err := task(); err != nil {
			if _, ok := <-errorChan; ok {
				errorChan <- err
			}
		}
	}
}
