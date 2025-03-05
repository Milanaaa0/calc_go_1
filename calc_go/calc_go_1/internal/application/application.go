package application

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	calculate "github.com/pashapdev/calc_go/pkg/calculation"
)

type Request struct {
	Expression string `json:"expression"`
}

type Response struct {
	ID string `json:"id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Expression struct {
	ID     string
	Status string
	Result float64
}

type Task struct {
	ID            string
	Arg1          float64
	Arg2          float64
	Operation     string
	OperationTime time.Duration
}

var (
	expressions = make(map[string]Expression)
	tasks       = make(map[string]calculate.Task)
	results     = make(map[string]calculate.Result)
	mutex       = &sync.Mutex{}
)

func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id := generateID()
	mutex.Lock()
	expressions[id] = Expression{ID: id, Status: "pending"}
	mutex.Unlock()

	go func() {
		result, err := calculate.Calc(req.Expression)
		mutex.Lock()
		defer mutex.Unlock()
		if err != nil {
			expressions[id] = Expression{ID: id, Status: "error", Result: 0}
		} else {
			expressions[id] = Expression{ID: id, Status: "done", Result: result}
		}
	}()

	resp := Response{ID: id}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func GetExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	var exprs []Expression
	for _, expr := range expressions {
		exprs = append(exprs, expr)
	}

	resp := struct {
		Expressions []Expression `json:"expressions"`
	}{Expressions: exprs}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func GetExpressionByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/v1/expressions/"):]
	mutex.Lock()
	defer mutex.Unlock()

	expr, exists := expressions[id]
	if !exists {
		http.Error(w, "Expression not found", http.StatusNotFound)
		return
	}

	resp := struct {
		Expression Expression `json:"expression"`
	}{Expression: expr}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	for _, task := range tasks {
		resp := struct {
			Task calculate.Task `json:"task"`
		}{Task: task}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}

	http.Error(w, "No tasks available", http.StatusNotFound)
}

func PostResultHandler(w http.ResponseWriter, r *http.Request) {
	var result calculate.Result
	err := json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	results[result.ID] = result
	w.WriteHeader(http.StatusOK)
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
