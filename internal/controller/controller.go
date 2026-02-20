package controller

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/httpListener"
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/scanners/blockURL"
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/scanners/typosquatting"

)

type Controller struct {
    Scans []Scan
}

type Scan interface {
    Scan(r *http.Request) (res bool, reasons []string)
}


func New() *Controller {
    scans := make([]Scan, 0, 4)
    scans = append(scans, blockURL.New())	//TODO: inject database service hear
    scans = append(scans, typosquatting.New())	//TODO: inject database service hear
    return &Controller{Scans: scans}
}

func (c *Controller) DisplayOperations() {
	reader := bufio.NewReader(os.Stdin)

	var shutdown func(context.Context) error
	serverRunning := false

	for {
		fmt.Printf(
			"<1> Passive Scan of Network\n" +
				"<2> Write block URL's\n" +
				"<3> Read blocked URL's\n" +
				"<4> Get History of Visited Domain\n" +
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
			shutdown, err = httplistener.ScanHTTPNetwork(c)
			if err != nil {
				fmt.Println("Failed to start server:", err)
				continue
			}
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

// InspectGET implements httplistener.Inspector.
func (c *Controller) InspectGET(req *http.Request) (res bool, reason string) {
   	var reasons []string
	res = true
    for _, s := range c.Scans {
        block, rs := s.Scan(req)
        reasons = append(reasons, rs...)
		if !block{
			res = false
		}
    }
    return res, parseMsg(reasons)
}

// InspectPOST implements httplistener.Inspector.
func (c *Controller) InspectPOST(req *http.Request) (res bool, reason string) {
	panic("unimplemented")
}

func parseMsg(reasons []string) (msg string) {
	panic("not implemented")
}
