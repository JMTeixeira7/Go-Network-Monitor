package httplistener

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/url"
	urlerr "github.com/JMTeixeira7/Go-Network-Monitor.git/internal/url_error"
)

type ProxyHandler struct {
    Inspector Inspector
}

const max_t = 10

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

	var body io.Reader	//ignores body in GET method even if it exists
	webRequest, err := http.NewRequest(req.Method, req.URL.String(), body)
	if err != nil {
		fmt.Printf("Could not create a new request: %s\n", err)
		return
	}
	webRequest.Header = req.Header.Clone()

	//TODO: check if visited in the last minute
		//send to controller for scan
	
	

	var webRes *http.Response
	webRes, err = sendRequest(webRequest)
	if err != nil {
		fmt.Printf("Error while redirecting request: %s\n", err)
		return
	}

	//TODO: assyncrounly add to the recently visited (last minute)

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

	//TODO: phishing + XSS scan

	var webRes *http.Response
	webRes, err = sendRequest(webRequest)
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

func sendRequest(webReq *http.Request) (*http.Response, error){
	client := &http.Client{Timeout: time.Duration(500) * time.Millisecond}
	var webRes *http.Response
	var err error
	for i := 0 ; i < max_t ; i++ { 
		webRes, err = client.Do(webReq)
		client.Timeout = client.Timeout*2
		if err != nil {
			if isTimeoutError(err){
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

func searchTargetURL(ctx context.Context, req *http.Request) error {
	url_target, err := url.CreateUrl(req.URL.String())
	if err != nil { 
		fmt.Printf("%s: Could not parse Url correctly: %s\n", ctx.Value(KeyServerAddr), err)
		return err
	}
	for _, target := range url.Urls {
		if !target.Target {
			continue
		}
		if url_target.Domain == target.Domain {
			fmt.Printf("%s: Requested Url is targetted: %s\n", ctx.Value(KeyServerAddr), url_target.Domain)
			return urlerr.NewTargetUrlError(req.RequestURI, target.Domain)
		}
	}
	return err
}