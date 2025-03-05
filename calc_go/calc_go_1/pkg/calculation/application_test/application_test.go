package application_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	application "github.com/pashapdev/calc_go/internal/application"
)

type Expression struct {
	ID     string
	Status string
	Result float64
}

func generateID() string {
	return "test-id"
}

var expressions = make(map[string]Expression)

func TestCalculateHandler(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  map[string]string
		expectedCode int
	}{
		{"valid expression", map[string]string{"expression": "2 + 2"}, http.StatusCreated},
		{"invalid expression", map[string]string{"expression": "2 + "}, http.StatusUnprocessableEntity},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req, err := http.NewRequest("POST", "/calculate", bytes.NewBuffer(body))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(application.CalculateHandler)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedCode)
			}
		})
	}
}

func TestGetExpressionsHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/expressions", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(application.GetExpressionsHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestGetExpressionByIDHandler(t *testing.T) {

	id := generateID()
	expressions[id] = Expression{ID: id, Status: "done", Result: 4}

	req, err := http.NewRequest("GET", "/expressions/"+id, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(application.GetExpressionByIDHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
