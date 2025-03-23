package main

import (
	"log"
	"net/http"

	"github.com/pashapdev/calc_go/internal/orchestrator"
)

func main() {

	orc := orchestrator.NewOrchestrator()
	handler := orchestrator.NewHandler(orc)

	http.HandleFunc("/task", handler.AddTask)
	http.HandleFunc("/internal/task", handler.GetTask)
	http.HandleFunc("/internal/status", handler.UpdateTaskStatus)
	http.HandleFunc("/internal/result", handler.AddResult)
	http.HandleFunc("/result", handler.GetResult)

	log.Println("Orchestrator started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
