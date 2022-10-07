package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	errorCh := make(chan struct{}, m)
	defer close(errorCh)
	wg := sync.WaitGroup{}
	defer wg.Wait()
	taskCh := make(chan Task)
	defer close(taskCh)
	wg.Add(n)

	for i := 1; i <= n; i++ {
		go doTask(&wg, taskCh, errorCh)
	}

	errCount := 0
	for _, task := range tasks {
		select {
		case <-errorCh:
			errCount++
			if m > 0 && errCount >= m {
				return ErrErrorsLimitExceeded
			}
		default:
		}
		taskCh <- task
	}

	return nil
}

func doTask(wg *sync.WaitGroup, taskCh chan Task, errorCh chan struct{}) {
	defer wg.Done()
	for task := range taskCh {
		taskError := task()
		if taskError != nil {
			select {
			case errorCh <- struct{}{}:
			default:
			}
		}
	}
}
