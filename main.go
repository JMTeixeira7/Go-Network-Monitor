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


	http.HandleFunc("/", SimpleRequestHandler)
	mux := http.NewServeMux()
	mux.HandleFunc("/", SimpleRequestHandler)
	mux.HandleFunc("/alternative", AlternativeRequestHandler)

	ctx, cancelCtx := context.WithCancel(context.Background())

	proxyServer := &http.Server{
		Addr: "127.0.0.1:3333",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, KeyServerAddr, l.Addr().String())
			return ctx
		},
	}

	err = proxyServer.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed){
		fmt.Println("Server is closed")
		} else if err!=nil {
			fmt.Println("Error while starting server: %s\n", err)
	}
	cancelCtx()
}
