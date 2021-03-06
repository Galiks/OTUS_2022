package hw05parallelexecution

import (
	"errors"
	"sync"
)

var (
	ErrErrorsLimitExceeded     = errors.New("errors limit exceeded")
	ErrNegativeGoroutineNumber = errors.New("goroutines' number is negative")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n < 0 {
		return ErrNegativeGoroutineNumber
	} else if n == 0 {
		return nil
	}
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}
	var (
		taskChan  = make(chan Task, len(tasks))
		errorChan = make(chan error, m+n)
		wg        sync.WaitGroup
	)
	defer close(errorChan)
	go func() {
		for _, task := range tasks {
			taskChan <- task
		}
		close(taskChan)
	}()

	for i := 0; i < n; i++ {
		wg.Add(1)
		go handler(taskChan, errorChan, m, &wg)
	}

	wg.Wait()
	if len(errorChan) >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func handler(taskChan chan Task, errorChan chan error, maxErrorCount int, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range taskChan {
		if err := task(); err != nil {
			errorChan <- err
		}
		if len(errorChan) >= maxErrorCount {
			return
		}
	}
}
