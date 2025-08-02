package main

import (
	//"fmt"
	"applications_manager/Url"
	"errors"
	"fmt"
	"net/http"
	"context"
	"net"
)

func main() {

	storage := Url.FileStorage{Filename: "data/urls.json"}

	Url.SetTargetURLs()
	err := storage.Save(Url.GetTargetURLs())
	if err != nil {
		fmt.Println("Failed to save URLs: ", err)
	}

	urls, err := storage.Load()
	if err != nil {
		fmt.Println("Failed to load URLs: ", err)
	}
	fmt.Println("Loaded URLs:", urls)

	//Multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/", DefaultRequestHandler)
	mux.HandleFunc("/Form", FormRequestHandler)
	mux.HandleFunc("/Body", ReadBodyRequestHandler)
	mux.HandleFunc("/QueryStrings", QueryStringRequestHandler)

	ctx, cancelCtx := context.WithCancel(context.Background())

	//httpServer
	proxyServer := &http.Server{
		Addr: "127.0.0.1:4444",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, KeyServerAddr, l.Addr().String())
			return ctx
		},
	}

	go func() {
		err = proxyServer.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed){
			fmt.Println("Server is closed")
			} else if err!=nil {
				fmt.Printf("Error while starting server: %s\n", err)
		}
		cancelCtx()
	}()

	<- ctx.Done() //Just to keep the program alive whilhe no behaviour exists
}
