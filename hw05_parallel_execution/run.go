package hw05parallelexecution

import (
	"errors"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	countOfTasks := 0
	countOfErrors := 0

	chEndTask := make(chan struct{}, n)
	chAddError := make(chan struct{}, n)

	lenTasks := len(tasks)
	finish := lenTasks
	if n < finish {
		finish = n
	}

	i := 0
	for ; i < finish; i++ {
		countOfTasks++
		go func(task Task) {
			defer func() { chEndTask <- struct{}{} }()
			if err := task(); err != nil {
				chAddError <- struct{}{}
			}
		}(tasks[i])
	}

	for {
		select {
		case <-chAddError:
			countOfErrors++
			if countOfErrors >= m {
				return ErrErrorsLimitExceeded
			}
		case <-chEndTask:
			countOfTasks--
			if i >= lenTasks {
				return nil
			}
			if countOfTasks < n {
				countOfTasks++
				go func() {
					defer func() { i++; chEndTask <- struct{}{} }()
					if err := tasks[i](); err != nil {
						chAddError <- struct{}{}
					}
				}()
			}
		}
	}
	// return nil
}
