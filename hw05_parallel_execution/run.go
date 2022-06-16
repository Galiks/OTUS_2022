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
		errorChan = make(chan<- error, m+n)
		wg        sync.WaitGroup
		// Или лучше использовать канал?
		globalErr error
		lock      sync.Mutex
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
		go handler(taskChan, errorChan, &globalErr, m, &wg, &lock)
	}

	wg.Wait()
	return globalErr
}

func handler(taskChan chan Task, errorChan chan<- error, globalErr *error, maxErrorCount int,
	wg *sync.WaitGroup, lock *sync.Mutex,
) {
	defer wg.Done()
	for task := range taskChan {
		lock.Lock()
		if *globalErr != nil {
			return
		}
		lock.Unlock()
		if err := task(); err != nil {
			errorChan <- err
		}
		if len(errorChan) >= maxErrorCount {
			lock.Lock()
			*globalErr = ErrErrorsLimitExceeded
			lock.Unlock()
			return
		}
	}
}
