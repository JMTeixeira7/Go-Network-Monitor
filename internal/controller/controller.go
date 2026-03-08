package controller

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/db/databaseService/blockUrlDBService"
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/db/databaseService/phishingDBService"
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/db/databaseService/typosquattingDBService"
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/model"
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/httpListener"
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/scanners/blockURL"
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/scanners/phishingPrevention"
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/scanners/typosquatting"
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/scanners/xssPrevention"
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/services/blockUrlAction"
)

type Controller struct {
	Scans   []Scan
	Actions map[string]ActionGroup
}

// InspectRequest implements httplistener.Inspector.
func (c *Controller) InspectRequest(req *http.Request) (res bool, reason string) {
	panic("unimplemented")
}

type Scan interface {
	Scan(r *http.Request) (res bool, reasons []string)
}

type ActionGroup interface {
	Name() string
}

type BlockActionUrlService interface {
	ActionGroup
	BlockUrl(domain string, schedules []model.Schedule) error
	GetAllBlockedURL() ([]string, error)
	GetBlockedURL(domain string) ([]model.Schedule, error)
}

func New(db *sql.DB) *Controller {
	scans := make([]Scan, 0, 4)
	scans = append(scans, xssPrevention.New(),
		blockURL.New(blockUrlDBService.NewBlockedDomainsDBService(db)),
		typosquatting.New(typosquattingDBService.NewTypoSquattingDBService(db)),
		phishingPrevention.New(phishingDBService.NewPhishingDBService(db)))
	actions := map[string]ActionGroup{
		"block_url_action": blockUrlAction.New(blockUrlDBService.NewBlockActionDomainsDBService(db)),
	}

	return &Controller{Scans: scans, Actions: actions}
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
				fmt.Println("Failed to start server: %s", err)
				continue
			}
			serverRunning = true
			fmt.Println("Server started on 127.0.0.1:4444")
		case "2":
			fmt.Printf("Introduce a domain you which to block with this format: e.g. google.com\n")
			line, err := reader.ReadString('\n')
			if err != nil {
				return
			}
			fmt.Printf("Introduce a scehdule o block domain: %s\n", line)
			fmt.Printf("Skip this if you don't wish to set a schedule\n\n")
			//TODO set schedule
			group, ok := c.Actions["block_url_action"]
			if !ok {
				fmt.Println("Did not find the Service for the given request")
				continue
			}
			blockGroup, ok := group.(BlockActionUrlService)
			if !ok {
				fmt.Println("Did not find the Service for the given request")
				continue
			}
			err = blockGroup.BlockUrl(line, nil)
			if err != nil {
				fmt.Println("Could not perform you request:\n%s\n", err)
				continue
			}
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
		if !block {
			res = false
		}
	}
	return res, parseMsg(reasons)
}

func parseMsg(reasons []string) (msg string) {
	panic("not implemented")
}
