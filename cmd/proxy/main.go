package main

import (
	"log"

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
    ctrl.DisplayOperations() // blocks here
}
