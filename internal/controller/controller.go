package controller

import (
	"fmt"
	"bufio"
	"context"

	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/httpListener"
)

func displayOperations() {
    reader := bufio.NewReader(os.Stdin)

    var shutdown func(context.Context) error
    serverRunning := false

    for {
        fmt.Printf(
            "<1> Passive Scan of Network\n"+
                "<2> Write block URL's\n"+
                "<3> Read blocked URL's\n"+
                "<4> Get History of Visited Domain\n"+
                "<5> Stop HTTP Server\n",
        )

        line, err := reader.ReadString('\n')
        if err != nil {
            return
        }

        choice := strings.TrimSpace(line)

        switch choice {
        case "1":
            if serverRunning {
                fmt.Println("Server already running.")
                continue
            }
            // ctrl implements httplistener.Inspector
            shutdown, _ = httplistener.scanHttpNetwork(ctrl)
            serverRunning = true
            fmt.Println("Server started on 127.0.0.1:4444")

        case "5":
            if !serverRunning {
                fmt.Println("Server not running.")
                continue
            }
            ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
            _ = shutdown(ctx)
            cancel()
            serverRunning = false
            fmt.Println("Server stopped.")

        default:
            fmt.Println("Not implemented yet.")
        }
    }
}
