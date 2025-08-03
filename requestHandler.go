package main

import (
	"fmt"
	"io"
	"net/http"
	"applications_manager/Url"
	"time"
	"errors"
	"net"
)


const KeyServerAddr = "ServerAdress"
const max_t = 10

func DefaultRequestHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	url, err := Url.CreateUrl(req.URL.String())
	if err != nil { 
		fmt.Printf("%s: Could not parse Url correctly: %s\n", ctx.Value(KeyServerAddr), err)
		return
	}

	for _, target := range Url.Urls {
		if !target.Target {
			break
		}
		if url.Domain == target.Domain {
			fmt.Printf("%s: Requested Url is targetted: %s\n", ctx.Value(KeyServerAddr), url.Domain)
		}
	}

	webRequest, err := http.NewRequest(req.Method, req.URL.String(), req.Body)
	if err != nil {
		fmt.Printf("Could not create a new request: %s\n, err")
		return
	}
	webRequest.Header = req.Header.Clone()


	client := &http.Client{Timeout: time.Duration(1) * time.Second}		//TODO: wrap in func
	var webRes *http.Response
	for i := 0 ; i < max_t + 1 ; i++ { 
		client.Timeout = client.Timeout*2
		webRes, err = client.Do(webRequest)
		if err != nil {
			if isTimeoutError(err){
				continue
			} else {
				return
			}
		} else if i == max_t {
			return
		}
		break
	}


	for key, value := range webRes.Header {		//TODO: wrap in func
		for _, b := range value {
			w.Header().Add(key, b)
		}
	}

	defer webRes.Body.Close()

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
