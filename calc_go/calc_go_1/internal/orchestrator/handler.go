package orchestrator

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/pashapdev/calc_go/internal/storage"
	"github.com/pashapdev/calc_go/internal/user"
	"github.com/pashapdev/calc_go/pkg/calculation"
)

type Handler struct {
	DB *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{DB: db}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if req.Login == "" || req.Password == "" {
		http.Error(w, "login and password required", http.StatusBadRequest)
		return
	}
	if err := user.Register(h.DB, req.Login, req.Password); err != nil {
		http.Error(w, "user exists or error", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if req.Login == "" || req.Password == "" {
		http.Error(w, "login and password required", http.StatusBadRequest)
		return
	}
	token, err := user.Login(h.DB, req.Login, req.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		claims, err := user.ParseJWT(parts[1])
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		next(w, r.WithContext(ctx))
	}
}

func (h *Handler) Calculate(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	var req struct {
		Expression string `json:"expression"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	result, err := calculation.Calc(req.Expression)
	if err != nil {
		http.Error(w, "calculation error", http.StatusBadRequest)
		return
	}

	storage.SaveCalculation(h.DB, userID, req.Expression, formatResult(result))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"result": result,
	})
}

func (h *Handler) History(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	history, err := storage.GetCalculationsByUser(h.DB, userID)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(history)
}

func formatResult(result float64) string {
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.8f", result), "0"), ".")
}
