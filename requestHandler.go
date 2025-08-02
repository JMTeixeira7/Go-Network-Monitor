package main

import (
	"fmt"
	"io"
	"net/http"
	"applications_manager/Url"
)


const KeyServerAddr = "ServerAdress"

func DefaultRequestHandler(res http.ResponseWriter, req *http.Request) {
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

	//perform DNS lookup

	//redirect message
	
	io.WriteString(res, fmt.Sprintf("Server: %s Received your Request\n", ctx.Value(KeyServerAddr)))
}


func FormRequestHandler(res http.ResponseWriter, req *http.Request) {
	//read form data
	myInstitution := req.PostFormValue("myInstitution")
	if myInstitution == "" {
		res.Header().Set("x-missing-field", "myInstitution")
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	io.WriteString(res, fmt.Sprintf("This is your institution, %s\n", myInstitution))
}

func QueryStringRequestHandler(res http.ResponseWriter, req *http.Request) {
	//get the URL query strings
	hasName := req.URL.Query().Has("name")
	name := req.URL.Query().Get("name")
	hasNumber := req.URL.Query().Has("number")
	number := req.URL.Query().Get("number")
	
	io.WriteString(res, fmt.Sprintf("QueryString has:\n Name(%t): %s\n Number(%t): %s\n", hasName, name, hasNumber, number))
}

func ReadBodyRequestHandler(res http.ResponseWriter, req *http.Request) {
	//read request body
	body, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("Server could not read the body: %s\n", err)
	}
	
	io.WriteString(res, fmt.Sprintf("This is the body of the Request, %s\n", body))
}