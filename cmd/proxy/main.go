package main

import (
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/controller"
)

func main() {
    ctrl := controller.New()
    ctrl.DisplayOperations() // blocks here
}
