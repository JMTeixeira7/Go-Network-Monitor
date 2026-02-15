package main

import (
	"errors"
	"fmt"
	"net/http"
	"context"
	"net"

	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/controller"
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/httpListener"
)

func main() {
	ctx, cancelCtx := context.WithCancel(context.Background())
	controller.displayOperations()
	<- ctx.Done() //Just to keep the program alive whilhe no behaviour exists
}
