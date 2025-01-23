package hw05parallelexecution

import (
	"errors"
	"sync"
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
	doneCh := make(chan error, n)
	stopTaskCh := make(chan struct{})
	wg := sync.WaitGroup{}
	defer func(w *sync.WaitGroup) {
		w.Wait()
		close(stopTaskCh)
		close(doneCh)
	}(&wg)

	for i := 0; i < n; i++ {
		go worker(taskCh, doneCh, &wg)
	}

	go func(taskCh1 chan Task, stopTaskCh1 chan struct{}) {
		defer close(taskCh1)
		for _, task := range tasks {
			select {
			case <-stopTaskCh1:
				return
			default:
				taskCh1 <- task
			}
		}
	}(taskCh, stopTaskCh)

	errCount := 0
	for err := range doneCh {
		if m > 0 {
			if err != nil {
				errCount++
				if errCount >= m {
					stopTaskCh <- struct{}{}
					return ErrErrorsLimitExceeded
				}
			}
		}
	}

	return nil
}

func worker(taskCh chan Task, doneCh chan error, w *sync.WaitGroup) {
	for task := range taskCh {
		w.Add(1)
		doneCh <- task()
		w.Done()
	}
}
