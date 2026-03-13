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

type Scan interface {
	Scan(r *http.Request) (res bool, reasons []string)
}

type ActionGroup interface {
	Name() string
}

type BlockActionUrlService interface {
	ActionGroup
	BlockUrl(ctx context.Context, domain string, schedules []string) error
	GetAllBlockedURL(ctx context.Context) ([]string, error)
	GetBlockedURL(ctx context.Context, domain string) ([]string, error)
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
			fmt.Print("Introduce a domain you which to block with this format: e.g. google.com\n")
			line, err := reader.ReadString('\n')
			if err != nil {
				return
			}
			fmt.Printf("Do you wish to set a block schedule for this domain, %s? [Yes/No]\n", line)
			response, err := readBinaryResponse(reader)
			if err != nil {
				fmt.Print(err)
				continue
			}
			schedules := []string{}
			if response {
				schedules, err = readSchedule(reader)
			}
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
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			err = blockGroup.BlockUrl(ctx, line, schedules)
			if err != nil {
				fmt.Printf("Could not perform your request:\n%s\n", err)
				continue
			}
		case "3":
			fmt.Println("Enter a domain or skip to view all blocked domains:")
			line, err := reader.ReadString('\n')
			if err != nil {
				return
			}
		
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
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			if line != "\n" {
				schedules, err := blockGroup.GetBlockedURL(ctx, line)
				if err != nil {
					fmt.Printf("Could not perform your request:\n%s\n", err)
					continue
				}
				fmt.Printf(displaySchedules(schedules))
			} else {
				blocked_domains, err := blockGroup.GetAllBlockedURL(ctx)
				if err != nil {
					fmt.Printf("Could not perform your request:\n%s\n", err)
					continue
				}
				fmt.Printf(displayBlockedDomains(blocked_domains))
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
func (c *Controller) InspectRequest(req *http.Request) (res bool, reason string) {
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

func parseMsg(reasons []string) string {
	if len(reasons) == 0 {
		return "No scan reasons.\n"
	}

	var b strings.Builder
	b.WriteString("Scan results:\n")

	count := 0
	for _, reason := range reasons {
		reason = strings.TrimSpace(reason)
		if reason == "" {
			continue
		}
		count++
		b.WriteString(fmt.Sprintf("  %d. %s\n", count, reason))
	}

	if count == 0 {
		return "No scan reasons.\n"
	}

	return b.String()
}

func readBinaryResponse(reader *bufio.Reader) (bool, error) {
	for true {
		response, err := reader.ReadString('\n')
		if err != nil {
			return false, fmt.Errorf("Failed to read user response, %w", err) 
		}
		response = strings.ToLower(strings.TrimSpace(response))
		if response == "y" || response == "yes" {
			return true, nil
		}
		if response == "n" || response == "no" {
			return false, nil
		}
		fmt.Print("Please enter [Yes/No]:")
	}
	return false, nil
}

func readSchedule(reader *bufio.Reader) ([]string, error) {
	var schedules []string
	for {
		fmt.Print("Enter schedule as: <timestamp> <timestamp> <weekday>\nUse - to skip a field. Press Enter on an empty line to finish:\n")

		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("failed to read user response: %w", err)
		}

		line = strings.TrimSpace(line)
		if line == "" {
			return schedules, nil
		}
		schedules = append(schedules, line)
	}
}

func displaySchedules(schedules []string) string {
	if len(schedules) == 0 {
		return "Schedules: none\n"
	}

	var b strings.Builder
	b.WriteString("Schedules:\n")
	for i, schedule := range schedules {
		b.WriteString(fmt.Sprintf("  %d. %s\n", i+1, schedule))
	}
	return b.String()
}

func displayBlockedDomains(domains []string) string {
	if len(domains) == 0 {
		return "Blocked domains: none\n"
	}

	var b strings.Builder
	b.WriteString("Blocked domains:\n")
	for i, domain := range domains {
		b.WriteString(fmt.Sprintf("  %d. %s\n", i+1, domain))
	}
	return b.String()
}