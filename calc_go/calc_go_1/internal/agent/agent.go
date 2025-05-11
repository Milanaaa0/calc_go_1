package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/pashapdev/calc_go/pkg/calculation"
)

type Agent struct {
	client          *http.Client
	orchestratorURL string
}

func NewAgent(orchestratorURL string) *Agent {
	return &Agent{
		client:          &http.Client{Timeout: 5 * time.Second},
		orchestratorURL: orchestratorURL,
	}
}

func (a *Agent) Start() {
	for {
		task, err := a.fetchTask()
		if err != nil {
			log.Printf("Error fetching task: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		result, err := calculation.Calc(fmt.Sprintf("%f%s%f", task.Arg1, task.Operation, task.Arg2))
		if err != nil {
			log.Printf("Error executing task %s: %v", task.ID, err)
			continue
		}

		if err := a.sendResult(task.ID, result); err != nil {
			log.Printf("Error sending result: %v", err)
		}
	}
}

func (a *Agent) fetchTask() (calculation.Task, error) {
	resp, err := a.client.Get(a.orchestratorURL + "/internal/task")
	if err != nil {
		return calculation.Task{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return calculation.Task{}, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	var task calculation.Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return calculation.Task{}, err
	}

	return task, nil
}

func (a *Agent) sendResult(taskID string, result float64) error {
	payload := map[string]interface{}{
		"id":     taskID,
		"result": result,
	}

	jsonData, _ := json.Marshal(payload)
	resp, err := a.client.Post(
		a.orchestratorURL+"/internal/result",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
