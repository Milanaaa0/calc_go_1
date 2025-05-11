package main

import (
	"log"
	"net/http"

	"github.com/pashapdev/calc_go/internal/orchestrator"
	"github.com/pashapdev/calc_go/internal/storage"
)

func main() {
	db, err := storage.InitDB("calc.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	handler := orchestrator.NewHandler(db)

	http.HandleFunc("/api/v1/register", handler.Register)
	http.HandleFunc("/api/v1/login", handler.Login)
	http.HandleFunc("/api/v1/calculate", orchestrator.AuthMiddleware(handler.Calculate))
	http.HandleFunc("/api/v1/history", orchestrator.AuthMiddleware(handler.History))

	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}
