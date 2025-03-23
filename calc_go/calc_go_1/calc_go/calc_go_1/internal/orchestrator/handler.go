package orchestrator

import (
	"encoding/json"
	"net/http"

	"github.com/pashapdev/calc_go/pkg/calculation"
)

type Handler struct {
	orchestrator *Orchestrator
}

func NewHandler(orchestrator *Orchestrator) *Handler {
	return &Handler{
		orchestrator: orchestrator,
	}
}

func (h *Handler) AddTask(w http.ResponseWriter, r *http.Request) {
	var task calculation.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.orchestrator.AddTask(task)
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	task, err := h.orchestrator.GetTask()
	if err != nil {
		http.Error(w, "No tasks available", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *Handler) UpdateTaskStatus(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TaskID string `json:"task_id"`
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.orchestrator.UpdateTaskStatus(req.TaskID, req.Status); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) AddResult(w http.ResponseWriter, r *http.Request) {
	var result calculation.Result
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.orchestrator.AddResult(result)
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetResult(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("task_id")
	result, err := h.orchestrator.GetResult(taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(result)
}
