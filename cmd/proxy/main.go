package main

import (
	"log"
	"net/http"

	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/controller"
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/db"
	"github.com/joho/godotenv"
)

func main() {

    _ = godotenv.Load()
    db, err := storage.OpenMySQL()
    if err != nil { log.Fatal(err) }
    defer db.Close()
    ctrl := controller.New(db)
    api := controller.NewApi(ctrl)

	addr := "127.0.0.1:8081"
	log.Printf("REST API listening on http://%s", addr)

	if err := http.ListenAndServe(addr, api.Handler()); err != nil {
		log.Fatal(err)
	}
}
