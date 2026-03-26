package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Server struct {
	ctrl *Controller
	mux  *http.ServeMux
}

func NewApi(ctrl *Controller) *Server {
	s := &Server{
		ctrl: ctrl,
		mux:  http.NewServeMux(),
	}

	s.routes()
	return s
}

func (s *Server) routes() {
	s.mux.HandleFunc("/api/status", s.handleStatus)
}

func (s *Server) Handler() http.Handler {
	return withCORS(s.mux)
}

type StatusResponse struct {
	ListenerRunning bool   `json:"listenerRunning"`
	CacheStatus     string `json:"cacheStatus"`
	LastUpdated     string `json:"lastUpdated"`
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"error": "method not allowed",
		})
		return
	}

	resp := StatusResponse{
		ListenerRunning: false,
		CacheStatus:     "unknown",
		LastUpdated:     time.Now().UTC().Format(time.RFC3339),
	}

	writeJSON(w, http.StatusOK, resp)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("failed to encode json response: %v", err)
	}
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}