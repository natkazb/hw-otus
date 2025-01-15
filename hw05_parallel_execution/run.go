package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n int, m int) error {
	if len(tasks) == 0 {
		return nil
	}

	taskCh := make(chan Task)
	doneCh := make(chan error, n)
	wg := sync.WaitGroup{}
	defer func(w *sync.WaitGroup) {
		w.Wait()
		close(doneCh)
		close(taskCh)
	}(&wg)

	for i := 0; i < n; i++ {
		go worker(taskCh, doneCh, &wg)
	}

	errCount := 0
	tasksLeft := len(tasks)

	idxTask := 0
	for tasksLeft > 0 {
		select {
		case err := <-doneCh:
			tasksLeft--
			if err != nil {
				errCount++
				if errCount >= m {
					return ErrErrorsLimitExceeded
				}
			}
		default:
			if idxTask < len(tasks) {
				taskCh <- tasks[idxTask]
				idxTask++
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
