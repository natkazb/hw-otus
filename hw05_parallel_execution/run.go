package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

/*
	Run starts tasks in n goroutines and stops its work when receiving m errors from tasks
	m <= 0 это знак игнорировать ошибки в принципе
*/

func Run(tasks []Task, n int, m int) error {
	if len(tasks) == 0 {
		return nil
	}

	taskCh := make(chan Task)
	wg := sync.WaitGroup{}

	var errCount int64
	for i := 0; i < n; i++ {
		go worker(taskCh, &wg, m, &errCount)
	}

	for _, t := range tasks {
		if m > 0 {
			if atomic.LoadInt64(&errCount) >= int64(m) {
				break
			}
		}
		taskCh <- t
	}
	close(taskCh)

	wg.Wait()

	if m > 0 {
		if atomic.LoadInt64(&errCount) >= int64(m) {
			return ErrErrorsLimitExceeded
		}
	}

	return nil
}

func worker(taskCh chan Task, w *sync.WaitGroup, m int, errCount *int64) {
	w.Add(1)
	for task := range taskCh {
		err := task()
		if m > 0 {
			if err != nil {
				atomic.AddInt64(errCount, 1)
			}
		}
	}
	w.Done()
}
