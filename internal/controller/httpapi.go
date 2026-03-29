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
	s.mux.HandleFunc("/api/cache/clear", s.handleClearCache)
	s.mux.HandleFunc("/api/listener/start", s.handleProxyStartUp)
	s.mux.HandleFunc("/api/listener/stop", s.handleProxyShutdown)
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
	cacheStatus := "active"
	if s.ctrl.isCacheCleared() {
		cacheStatus = "cleared"
	}
	resp := StatusResponse{
		ListenerRunning: s.ctrl.isProxyRunning(),
		CacheStatus:     cacheStatus,
		LastUpdated:     time.Now().UTC().Format(time.RFC3339),
	}
	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleProxyStartUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"message": "method not allowed",
			"data":    nil,
		})
		return
	}

	err := s.ctrl.runProxy()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "failed to start proxy",
			"data":    nil,
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"message": "proxy started successfully",
		"data": map[string]any{
			"listenerRunning": s.ctrl.isProxyRunning(),
		},
	})
}

func (s *Server) handleClearCache(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"message": "method not allowed",
			"data":    nil,
		})
		return
	}

	err := s.ctrl.clearCache(nil)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "proxy not running or failed to initialize",
			"data":    nil,
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"message": "cache cleaned successfully",
		"data": map[string]any{
			"cache cleaned": s.ctrl.isCacheCleared(),
		},
	})
}

func (s *Server) handleProxyShutdown(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"message": "method not allowed",
			"data":    nil,
		})
		return
	}

	err := s.ctrl.shutdownProxy()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "failed to stop proxy",
			"data":    nil,
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"message": "proxy stoped successfully",
		"data": map[string]any{
			"listenerRunning": s.ctrl.isProxyRunning(),
		},
	})
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
