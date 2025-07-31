package main

import (
	"fmt"
	"io"
	"net/http"
)


const KeyServerAddr = "ServerAdress"

func SimpleRequestHandler(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	fmt.Println("Request received on Server: ", ctx.Value(KeyServerAddr))
	io.WriteString(res, "received on Handler1\n")
}


func AlternativeRequestHandler(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	fmt.Println("Request received on Server: ", ctx.Value(KeyServerAddr))
	io.WriteString(res, "received on Handler2\n")
}