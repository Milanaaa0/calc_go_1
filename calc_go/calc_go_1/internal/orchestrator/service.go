package orchestrator

import (
	"errors"
	"sync"

	"github.com/pashapdev/calc_go/pkg/calculation"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type Orchestrator struct {
	tasks   map[string]calculation.Task
	results map[string]calculation.Result
	mu      sync.RWMutex
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		tasks:   make(map[string]calculation.Task),
		results: make(map[string]calculation.Result),
	}
}

func (o *Orchestrator) AddTask(task calculation.Task) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.tasks[task.ID] = task
}

func (o *Orchestrator) GetTask() (calculation.Task, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	for _, task := range o.tasks {
		if task.Status == "pending" {
			task.Status = "processing"
			o.tasks[task.ID] = task
			return task, nil
		}
	}
	return calculation.Task{}, ErrTaskNotFound
}

func (o *Orchestrator) UpdateTaskStatus(taskID string, status string) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	task, exists := o.tasks[taskID]
	if !exists {
		return ErrTaskNotFound
	}
	task.Status = status
	o.tasks[taskID] = task
	return nil
}

func (o *Orchestrator) AddResult(result calculation.Result) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.results[result.ID] = result
}

func (o *Orchestrator) GetResult(taskID string) (calculation.Result, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	result, exists := o.results[taskID]
	if !exists {
		return calculation.Result{}, ErrTaskNotFound
	}
	return result, nil
}
