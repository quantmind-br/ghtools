package runner

import (
	"fmt"
	"sync"
)

type TaskResult struct {
	Name    string
	Success bool
	Message string
}

type ParallelRunner struct {
	MaxJobs int
	Results []TaskResult
	mu      sync.Mutex
}

func New(maxJobs int) *ParallelRunner {
	return &ParallelRunner{MaxJobs: maxJobs}
}

type Task struct {
	Name string
	Fn   func() (string, error)
}

func (r *ParallelRunner) Run(tasks []Task, onProgress func(done, total int)) []TaskResult {
	sem := make(chan struct{}, r.MaxJobs)
	var wg sync.WaitGroup
	done := 0
	total := len(tasks)

	for _, task := range tasks {
		wg.Add(1)
		sem <- struct{}{}
		go func(t Task) {
			defer wg.Done()
			defer func() { <-sem }()

			msg, err := t.Fn()
			result := TaskResult{Name: t.Name, Success: err == nil}
			if err != nil {
				result.Message = fmt.Sprintf("Failed: %s", err)
			} else {
				result.Message = msg
			}

			r.mu.Lock()
			r.Results = append(r.Results, result)
			done++
			if onProgress != nil {
				onProgress(done, total)
			}
			r.mu.Unlock()
		}(task)
	}

	wg.Wait()
	return r.Results
}
