package agent

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pashapdev/calc_go/pkg/calculation"
)

func TestFetchTask_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		task := calculation.Task{
			ID:        "task1",
			ExprID:    "expr1",
			Arg1:      10,
			Arg2:      5,
			Operation: "+",
			Status:    "pending",
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(task)
	}))
	defer ts.Close()

	agent := NewAgent(ts.URL)

	task, err := agent.fetchTask()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if task.ID != "task1" || task.Arg1 != 10 || task.Arg2 != 5 || task.Operation != "+" {
		t.Errorf("Expected task1 with args 10 and 5, got %v", task)
	}
}

func TestFetchTask_Error(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()
	agent := NewAgent(ts.URL)
	_, err := agent.fetchTask()
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := "status code: 500"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%v'", expectedError, err)
	}
}

func TestSendResult_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	agent := NewAgent(ts.URL)
	err := agent.sendResult("task1", 15)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestSendResult_Error(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()
	agent := NewAgent(ts.URL)

	err := agent.sendResult("task1", 15)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := "unexpected status code: 500"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%v'", expectedError, err)
	}
}

func TestStart(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/internal/task" {

			task := calculation.Task{
				ID:        "task1",
				ExprID:    "expr1",
				Arg1:      10,
				Arg2:      5,
				Operation: "+",
				Status:    "pending",
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(task)
		} else if r.URL.Path == "/internal/result" {

			var payload map[string]interface{}
			json.NewDecoder(r.Body).Decode(&payload)
			if payload["id"] != "task1" || payload["result"] != 15.0 {
				t.Errorf("Expected result 15 for task1, got %v", payload)
			}
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer ts.Close()
	agent := NewAgent(ts.URL)

	go agent.Start()
	time.Sleep(1 * time.Second)
}
