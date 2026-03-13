package httplistener

import (
	//"context"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"
)

type ProxyHandler struct {
	Inspector Inspector
	cache map[string]time.Time
	cacheTTL time.Duration
	max_tries int
	mu *sync.RWMutex
}

func (h *ProxyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		h.GetHandler(w, req)
	case http.MethodPost:
		h.PostHandler(w, req)
	default:
		http.Error(w, "Method Not Supported", http.StatusMethodNotAllowed)
	}
}

func (h *ProxyHandler) GetHandler(w http.ResponseWriter, req *http.Request) {
	//ctx := req.Context()

	var body io.Reader //ignores body in GET method even if it exists
	webRequest, err := http.NewRequest(req.Method, req.URL.String(), body)
	if err != nil {
		fmt.Printf("Could not create a new request: %s\n", err)
		return
	}
	webRequest.Header = req.Header.Clone()

	//TODO: check if visited in the last minute
	//send to controller for scan
	if !h.seenRecently(req.URL.String()) { //Do this via interface?
		res, docs := h.Inspector.InspectRequest(webRequest)
		if res {
			fmt.Printf("Scanning results:\n %s\n", docs)
			return
		}
	}

	var webRes *http.Response
	webRes, err = h.sendRequest(webRequest)
	if err != nil {
		fmt.Printf("Error while redirecting request: %s\n", err)
		return
	}

	h.markSeen(req.URL.String())

	defer webRes.Body.Close()
	for key, value := range webRes.Header {
		for _, b := range value {
			w.Header().Set(key, b)
		}
	}
	w.WriteHeader(webRes.StatusCode)
	io.Copy(w, webRes.Body)
}

func (h *ProxyHandler) PostHandler(w http.ResponseWriter, req *http.Request) {
	webRequest, err := http.NewRequest(req.Method, req.URL.String(), req.Body)
	if err != nil {
		fmt.Printf("Could not create new request: %s", err)
		return
	}
	webRequest.Header = req.Header.Clone()

	res, docs := h.Inspector.InspectRequest(webRequest)
	if res {
		fmt.Printf("Scanning results:\n %s\n", docs)
		return
	}

	var webRes *http.Response
	webRes, err = h.sendRequest(webRequest)
	if err != nil {
		fmt.Printf("Error while redirecting request: %s\n", err)
		return
	}

	defer webRes.Body.Close()
	for key, value := range webRes.Header {
		for _, b := range value {
			w.Header().Set(key, b)
		}
	}
	w.WriteHeader(webRes.StatusCode)
	io.Copy(w, webRes.Body)
}


func isTimeoutError(err error) bool {
	// Check if error is a net timeout
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}

func (h *ProxyHandler) sendRequest(webReq *http.Request) (*http.Response, error) {
	client := &http.Client{Timeout: time.Duration(500) * time.Millisecond}
	var webRes *http.Response
	var err error
	for i := 0; i < h.max_tries; i++ {
		webRes, err = client.Do(webReq)
		client.Timeout = client.Timeout * 2
		if err != nil {
			if isTimeoutError(err) {
				continue
			} else {
				return nil, err
			}
		} else {
			return webRes, nil
		}
	}

	if err != nil {
		return nil, err
	}

	return webRes, err
}

func (h *ProxyHandler) markSeen(domain string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.cache == nil {
		h.cache = make(map[string]time.Time)
	}
	h.cache[domain] = time.Now()
}

func (h *ProxyHandler) seenRecently(domain string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if h.cache == nil {
		return false
	}
	now := time.Now()
	lastTime, ok := h.cache[domain]
	if !ok {
		return false
	}
	if now.Sub(lastTime) > h.cacheTTL {
		return false
	}
	return true
}

func (h *ProxyHandler) cleanupExpired() {
	now := time.Now()

	h.mu.Lock()
	defer h.mu.Unlock()

	for domain, t := range h.cache {
		if now.Sub(t) > h.cacheTTL {
			delete(h.cache, domain)
		}
	}
}

func (h *ProxyHandler) startCleanup(ctx context.Context, timer time.Duration) {
	ticker := time.NewTicker(timer)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			h.cleanupExpired()
		case <-ctx.Done():
			return
		}
	}
}
