package main

import (
	"log"

	"github.com/pashapdev/calc_go/internal/agent"
)

func main() {

	orchestratorURL := "http://localhost:8080"

	agent := agent.NewAgent(orchestratorURL)
	log.Println("Agent started and waiting for tasks...")
	agent.Start()
}
