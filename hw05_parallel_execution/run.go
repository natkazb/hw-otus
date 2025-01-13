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
	errCh := make(chan error)
	doneCh := make(chan struct{})
	wg := sync.WaitGroup{}

	// Запуск n воркеров
	for i := 0; i < n; i++ {
		go worker(taskCh, errCh, doneCh, &wg)
	}

	// Счетчик ошибок
	errCount := 0
	tasksLeft := len(tasks)

	// Отправка первой партии задач
	go func() {
		for _, t := range tasks {
			taskCh <- t
		}
		close(taskCh)
	}()

	// Обработка результатов
	for tasksLeft > 0 {
		select {
		case err := <-errCh:
			if err != nil {
				errCount++
				if errCount >= m {
					wg.Wait()
					return ErrErrorsLimitExceeded
				}
			}
			tasksLeft--
		case <-doneCh:
			tasksLeft--
		}
	}

	wg.Wait()
	return nil
}

func worker(taskCh <-chan Task, errCh chan<- error, doneCh chan<- struct{}, w *sync.WaitGroup) {
	for task := range taskCh {
		w.Add(1)
		err := task()
		if err != nil {
			w.Done()
			errCh <- err
		} else {
			w.Done()
			doneCh <- struct{}{}
		}
	}
}
