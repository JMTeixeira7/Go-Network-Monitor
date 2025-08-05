package main

import (
	"applications_manager/Url"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)


const KeyServerAddr = "ServerAdress"
const max_t = 10

func DefaultRequestHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		GetHandler(w, req)
	case http.MethodPost:
		PostHandler(w, req)
	case http.MethodPut:
		PutHandler(w, req)
	case http.MethodDelete:
		DeleteHandler(w, req)
	default:
		http.Error(w, "Method Not Supported", http.StatusMethodNotAllowed)
	}
}

func GetHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	err := searchTargetURL(ctx, req)
	if err != nil {
		return
	}

	var body io.Reader	//ignores body in GET method even if it exists
	webRequest, err := http.NewRequest(req.Method, req.URL.String(), body)
	if err != nil {
		fmt.Printf("Could not create a new request: %s\n", err)
		return
	}
	webRequest.Header = req.Header.Clone()


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


func FormRequestHandler(w http.ResponseWriter, req *http.Request) {
	//read form data
	myInstitution := req.PostFormValue("myInstitution")
	if myInstitution == "" {
		w.Header().Set("x-missing-field", "myInstitution")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	io.WriteString(w, fmt.Sprintf("This is your institution, %s\n", myInstitution))
}

func QueryStringRequestHandler(w http.ResponseWriter, req *http.Request) {
	//get the URL query strings
	hasName := req.URL.Query().Has("name")
	name := req.URL.Query().Get("name")
	hasNumber := req.URL.Query().Has("number")
	number := req.URL.Query().Get("number")
	
	io.WriteString(w, fmt.Sprintf("QueryString has:\n Name(%t): %s\n Number(%t): %s\n", hasName, name, hasNumber, number))
}

func ReadBodyRequestHandler(w http.ResponseWriter, req *http.Request) {
	//read request body
	body, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("Server could not read the body: %s\n", err)
	}
	
	io.WriteString(w, fmt.Sprintf("This is the body of the Request, %s\n", body))
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
	url, err := Url.CreateUrl(req.URL.String())
	if err != nil { 
		fmt.Printf("%s: Could not parse Url correctly: %s\n", ctx.Value(KeyServerAddr), err)
		return err
	}
	for _, target := range Url.Urls {
		if !target.Target {
			break
		}
		if url.Domain == target.Domain {
			fmt.Printf("%s: Requested Url is targetted: %s\n", ctx.Value(KeyServerAddr), url.Domain)
		}
	}
	return err
}