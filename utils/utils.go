package utils

import (
	"fmt"
	"sync"
)

func RunParallelTasks(taskMap map[string]Task, taskFunc func(Task) error, maxParallel int) error {
	var wg sync.WaitGroup
	var errs []error
	sem := make(chan struct{}, maxParallel)

	for _, task := range taskMap {
		wg.Add(1)
		go func(t Task) {
			defer wg.Done()
			sem <- struct{}{}        // Acquire semaphore
			defer func() { <-sem }() // Release semaphore
			fmt.Printf("Running task %s (type: %s)\n", t.ID, t.Type)

			if err := taskFunc(t); err != nil {
				errs = append(errs, err)
			}
		}(task)
	}

	wg.Wait()

	if len(errs) > 0 {
		return fmt.Errorf("errors occurred: %v", errs)
	}

	return nil
}

func RunParallel(fns []func() error, maxParallel int) error {
	var wg sync.WaitGroup
	var errs []error
	sem := make(chan struct{}, maxParallel)

	for _, fn := range fns {
		wg.Add(1)
		go func(f func() error) {
			defer wg.Done()
			sem <- struct{}{}        // Acquire semaphore
			defer func() { <-sem }() // Release semaphore

			if err := f(); err != nil {
				errs = append(errs, err)
			}
		}(fn)
	}

	wg.Wait()

	if len(errs) > 0 {
		return fmt.Errorf("errors occurred: %v", errs)
	}

	return nil
}
