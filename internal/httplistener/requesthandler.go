package httplistener

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/proxyError"
)

type CacheCommand struct {
	DeleteDomains []string
	ClearAll      bool
}


type Handler struct {
	inspector      Inspector
	cache          map[string]time.Time
	cacheTTL       time.Duration
	maxRetries     int
	cacheCommands  chan CacheCommand
	errorResponder proxyError.ErrorResponder

	mu sync.RWMutex
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	responder := h.errorResponder
	if responder == nil {
		responder = proxyError.PlainTextErrorResponder{}
	}

	if err := h.handle(w, r); err != nil {
		log.Printf("proxy request failed: method=%s host=%s err=%v", r.Method, r.URL.Host, err)
		responder.WriteError(w, r, err)
	}
}

func (h *Handler) handle(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return h.handleGet(w, r)
	case http.MethodPost:
		return h.handlePost(w, r)
	default:
		return &proxyError.HTTPError{
			Status:  http.StatusMethodNotAllowed,
			Message: "method not allowed",
		}
	}
}

func (h *Handler) handleGet(w http.ResponseWriter, r *http.Request) error {
	outReq, err := newOutboundRequest(r, nil)
	if err != nil {
		return &proxyError.HTTPError{
			Status:  http.StatusBadRequest,
			Message: "invalid request",
			Err:     err,
		}
	}

	cacheKey := r.URL.Host
	if !h.wasSeenRecently(cacheKey) {
		allowed, docs := h.inspector.InspectRequest(outReq)
		if !allowed {
			return &proxyError.HTTPError{
				Status:  http.StatusForbidden,
				Message: "request blocked",
				Err:     fmt.Errorf("inspection rejected request: %s", docs),
			}
		}
		h.markSeen(cacheKey)
	}

	outRes, err := h.sendRequest(outReq)
	if err != nil {
		return &proxyError.HTTPError{
			Status:  http.StatusBadGateway,
			Message: "upstream request failed",
			Err:     err,
		}
	}
	defer outRes.Body.Close()

	return writeUpstreamResponse(w, outRes)
}

func (h *Handler) handlePost(w http.ResponseWriter, r *http.Request) error {
	outReq, err := newOutboundRequest(r, r.Body)
	if err != nil {
		return &proxyError.HTTPError{
			Status:  http.StatusBadRequest,
			Message: "invalid request",
			Err:     err,
		}
	}

	allowed, docs := h.inspector.InspectRequest(outReq)
	if !allowed {
		return &proxyError.HTTPError{
			Status:  http.StatusForbidden,
			Message: "request blocked",
			Err:     fmt.Errorf("inspection rejected request: %s", docs),
		}
	}

	outRes, err := h.sendRequest(outReq)
	if err != nil {
		return &proxyError.HTTPError{
			Status:  http.StatusBadGateway,
			Message: "upstream request failed",
			Err:     err,
		}
	}
	defer outRes.Body.Close()

	return writeUpstreamResponse(w, outRes)
}

func newOutboundRequest(in *http.Request, body io.Reader) (*http.Request, error) {
	outReq, err := http.NewRequestWithContext(in.Context(), in.Method, in.URL.String(), body)
	if err != nil {
		return nil, fmt.Errorf("build outbound request: %w", err)
	}

	outReq.Header = in.Header.Clone()
	return outReq, nil
}

func writeUpstreamResponse(w http.ResponseWriter, res *http.Response) error {
	copyHeaders(w.Header(), res.Header)
	w.WriteHeader(res.StatusCode)

	if _, err := io.Copy(w, res.Body); err != nil {
		return fmt.Errorf("copy upstream response body: %w", err)
	}

	return nil
}

func copyHeaders(dst, src http.Header) {
	for key, values := range src {
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}

func isTimeoutError(err error) bool {
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}


/*
	TODO: http is stateless, while retrying we should make sure the operation doesnt have double effect problem
*/
func (h *Handler) sendRequest(req *http.Request) (*http.Response, error) {
	retries := h.maxRetries
	if retries < 1 {
		retries = 1
	}

	client := &http.Client{}
	timeout := 500 * time.Millisecond
	// Safe retries are only guaranteed when the body can be recreated.
	canRetry := req.Body == nil || req.GetBody != nil

	for attempt := 1; attempt <= retries; attempt++ {
		var attemptReq *http.Request

		if attempt > 1 {
			if !canRetry {
				break
			}

			body, err := req.GetBody()
			if err != nil {
				return nil, fmt.Errorf("reset request body for retry: %w", err)
			}

			attemptReq = req.Clone(req.Context())
			attemptReq.Body = body
		}

		client.Timeout = timeout
		res, err := client.Do(attemptReq)
		if err == nil {
			return res, nil
		}
		if !isTimeoutError(err) {
			return nil, fmt.Errorf("send upstream request: %w", err)
		}
		if attempt >= retries {
			return nil, fmt.Errorf("send upstream request after %d attempts: %w", attempt, err)
		}
		timeout *= 2 //exponencial backoff
	}
	return nil, fmt.Errorf("request body is not replayable, retry skipped")
}


/*
	Cache Management: Includes multiple features and operation of the http cache
*/

func (h *Handler) markSeen(key string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.cache == nil {
		h.cache = make(map[string]time.Time)
	}
	h.cache[key] = time.Now()
}

func (h *Handler) wasSeenRecently(key string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if h.cache == nil {
		return false
	}

	lastSeen, ok := h.cache[key]
	if !ok {
		return false
	}

	return time.Since(lastSeen) <= h.cacheTTL
}

func (h *Handler) cleanupExpiredCache() {
	now := time.Now()

	h.mu.Lock()
	defer h.mu.Unlock()

	for key, t := range h.cache {
		if now.Sub(t) > h.cacheTTL {
			delete(h.cache, key)
		}
	}
}

func (h *Handler) applyCacheCommand(cmd CacheCommand) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	switch {
	case cmd.ClearAll:
		clear(h.cache)
		return nil

	case len(cmd.DeleteDomains) > 0:
		for _, domain := range cmd.DeleteDomains {
			delete(h.cache, domain)
		}
		return nil

	default:
		return fmt.Errorf("invalid cache command")
	}
}

func (h *Handler) startCacheRoutine(ctx context.Context, every time.Duration) {
	ticker := time.NewTicker(every)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			h.cleanupExpiredCache()

		case cmd := <-h.cacheCommands:
			if err := h.applyCacheCommand(cmd); err != nil {
				log.Printf("cache command failed: %v", err)
			}

		case <-ctx.Done():
			return
		}
	}
}
