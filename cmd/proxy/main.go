package main

import (
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/controller"
)

func main() {
    ctrl := controller.New(fp)
    ctrl.DisplayOperations() // blocks here
}
