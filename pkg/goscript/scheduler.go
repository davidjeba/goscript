package goscript

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// TaskHandler runs a unit of work.
type TaskHandler func(context.Context) (interface{}, error)

// Task describes a scheduled job.
type Task struct {
	ID       string
	Name     string
	Priority int
	Handler  TaskHandler
	Metadata map[string]string
}

// TaskResult stores the outcome of a task.
type TaskResult struct {
	TaskID    string
	Name      string
	Value     interface{}
	Err       error
	StartedAt time.Time
	FinishedAt time.Time
}

// Scheduler executes tasks across worker goroutines.
type Scheduler struct {
	workers int
	tasks   chan Task
	results chan TaskResult
	stop    chan struct{}
	once    sync.Once
	wg      sync.WaitGroup
}

// NewScheduler creates a scheduler with a worker count.
func NewScheduler(workers int) *Scheduler {
	if workers < 1 {
		workers = 1
	}

	return &Scheduler{
		workers: workers,
		tasks:   make(chan Task, workers*4),
		results: make(chan TaskResult, workers*4),
		stop:    make(chan struct{}),
	}
}

// Start begins worker processing.
func (s *Scheduler) Start(ctx context.Context) {
	for i := 0; i < s.workers; i++ {
		s.wg.Add(1)
		go s.worker(ctx)
	}
}

// Submit queues a task for execution.
func (s *Scheduler) Submit(task Task) error {
	if task.Handler == nil {
		return fmt.Errorf("task handler is required")
	}

	if task.ID == "" {
		task.ID = fmt.Sprintf("task-%d", time.Now().UnixNano())
	}

	select {
	case <-s.stop:
		return fmt.Errorf("scheduler is stopped")
	case s.tasks <- task:
		return nil
	}
}

// Results exposes task results.
func (s *Scheduler) Results() <-chan TaskResult {
	return s.results
}

// Stop shuts down the scheduler.
func (s *Scheduler) Stop() {
	s.once.Do(func() {
		close(s.stop)
	})
	s.wg.Wait()
	close(s.results)
}

func (s *Scheduler) worker(ctx context.Context) {
	defer s.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.stop:
			return
		case task := <-s.tasks:
			started := time.Now()
			value, err := task.Handler(ctx)
			result := TaskResult{
				TaskID:    task.ID,
				Name:      task.Name,
				Value:     value,
				Err:       err,
				StartedAt: started,
				FinishedAt: time.Now(),
			}

			select {
			case s.results <- result:
			case <-ctx.Done():
				return
			case <-s.stop:
				return
			}
		}
	}
}

