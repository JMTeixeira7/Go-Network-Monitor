package main

import (
	//"fmt"
	"errors"
	"fmt"
	"net/http"
	"context"
	"net"

	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/handler"
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/storage"
)

func main() {

	storage := storage.FileStorage{Filename: "data/urls.json"}

	urls, err := storage.LoadAll()
	if err != nil {
		fmt.Println("Failed to load URLs: ", err)
	}
	fmt.Println("Loaded URLs:", urls)

	//Multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.DefaultRequestHandler)
	mux.HandleFunc("/Form", handler.FormRequestHandler)
	mux.HandleFunc("/Body", handler.ReadBodyRequestHandler)
	mux.HandleFunc("/QueryStrings", handler.QueryStringRequestHandler)

	ctx, cancelCtx := context.WithCancel(context.Background())

	//httpServer
	proxyServer := &http.Server{
		Addr: "127.0.0.1:4444",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, handler.KeyServerAddr, l.Addr().String())
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
