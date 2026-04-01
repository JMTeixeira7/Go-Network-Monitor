package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/controller/dto"
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
	s.mux.HandleFunc("/api/blocked-domains/", s.handleFetchBlockedDomain)
	s.mux.HandleFunc("/api/blocked-domains", s.handleFetchBlockedDomains)
}

func (s *Server) Handler() http.Handler {
	return withCORS(s.mux)
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
	resp := dto.StatusResponse{
		ListenerRunning: s.ctrl.isProxyRunning(),
		CacheStatus:     cacheStatus,
		LastUpdated:     time.Now().UTC().Format(time.RFC3339),
	}
	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleProxyStartUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, dto.ApiResponse[*dto.ListenerStateResponse]{
			Success: false,
			Message: "method not allowed",
			Data:    nil,
		})
		return
	}

	err := s.ctrl.runProxy()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, dto.ApiResponse[*dto.ListenerStateResponse]{
			Success: false,
			Message: "failed to start proxy",
			Data:    nil,
		})
		return
	}

	writeJSON(w, http.StatusOK, dto.ApiResponse[*dto.ListenerStateResponse]{
		Success: true,
		Message: "proxy started successfully",
		Data: &dto.ListenerStateResponse{
			ListenerRunning: s.ctrl.isProxyRunning(),
		},
	})
}

func (s *Server) handleClearCache(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, dto.ApiResponse[*dto.CacheStateResponse]{
			Success: false,
			Message: "method not allowed",
			Data: &dto.CacheStateResponse{},
		})
		return
	}

	err := s.ctrl.clearCache(nil)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, dto.ApiResponse[*dto.CacheStateResponse]{
			Success: false,
			Message: "proxy not running or failed to initialize",
			Data: &dto.CacheStateResponse{},
		})
		return
	}

	writeJSON(w, http.StatusOK, dto.ApiResponse[*dto.CacheStateResponse]{
		Success: true,
		Message: "cache cleaned successfully",
		Data: &dto.CacheStateResponse{
			CacheCleared: s.ctrl.isCacheCleared(),
		},
	})
}

func (s *Server) handleProxyShutdown(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, dto.ApiResponse[*dto.ListenerStateResponse]{
			Success: false,
			Message: "method not allowed",
			Data: nil,
		})
		return
	}

	err := s.ctrl.shutdownProxy()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, dto.ApiResponse[dto.ListenerStateResponse]{
			Success: false,
			Message: "failed to stop proxy",
			Data:    dto.ListenerStateResponse{},
		})
		return
	}

	writeJSON(w, http.StatusOK, dto.ApiResponse[dto.ListenerStateResponse]{
		Success: true,
		Message: "proxy stopped successfully",
		Data: dto.ListenerStateResponse{
			ListenerRunning: s.ctrl.isProxyRunning(),
		},
	})
}

func (s *Server) handleFetchBlockedDomains(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var req dto.BlockedDomainRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, dto.ApiResponse[*dto.BlockedDomainResponse]{
				Success: false,
				Message: "invalid JSON body",
				Data:    nil,
			})
		} else {
			res, err := s.ctrl.blockDomain(req)
			if err != nil {
				writeJSON(w, http.StatusInternalServerError, dto.ApiResponse[*dto.BlockedDomainResponse]{
					Success: false,
					Message: "failed to fetch blocked domains",
					Data:    nil,
				})
			} else {
				writeJSON(w, http.StatusOK, dto.ApiResponse[*dto.BlockedDomainResponse]{
					Success: true,
					Message: "domain blocked successfully",
					Data:    res,
				})
			}
		}
		return
	}

	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, dto.ApiResponse[*dto.BlockedDomainResponse]{
			Success: false,
			Message: "method not allowed",
			Data: nil,
		})
		return
	}

	res, err := s.ctrl.fetchBlockedDomains()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, dto.ApiResponse[[]dto.BlockedDomainResponse]{
			Success: false,
			Message: "failed to fetch blocked domains",
			Data:    res,
		})
	} else {
		writeJSON(w, http.StatusOK, dto.ApiResponse[[]dto.BlockedDomainResponse]{
			Success: true,
			Message: "blocked domains fetched successfully",
			Data:    res,
		})
	}
}

func (s *Server) handleFetchBlockedDomain(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, dto.ApiResponse[*dto.BlockedDomainResponse]{
			Success: false,
			Message: "method not allowed",
			Data: nil,
		})
		return
	}
	prefix := "/api/blocked-domains/"
	rawDomain := strings.TrimPrefix(r.URL.Path, prefix)
	if rawDomain == "" {
		writeJSON(w, http.StatusBadRequest, dto.ApiResponse[*dto.BlockedDomainResponse]{
			Success: false,
			Message: ("bad request: no domain specified"),
			Data: nil,
		})
		return
	}
	domain, err := url.PathUnescape(rawDomain)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.ApiResponse[*dto.BlockedDomainResponse]{
			Success: false,
			Message: "invalid domain path",
			Data:    nil,
		})
		return
	}
	res, err := s.ctrl.fetchBlockedDomain(domain)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, dto.ApiResponse[*dto.BlockedDomainResponse]{
			Success: false,
			Message: "failed to fetch blocked domain",
			Data:    nil,
		})
	} else {
		writeJSON(w, http.StatusOK, dto.ApiResponse[*dto.BlockedDomainResponse]{
			Success: true,
			Message: "blocked domain fetched successfully",
			Data:    res,
		})
	}
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
