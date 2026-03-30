package controller

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/controller/dto"
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/db/databaseService/blockUrlDBService"
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/db/databaseService/phishingDBService"
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/db/databaseService/visitDBService"
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/httplistener"
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/scanners/blockURL"
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/scanners/phishingPrevention"
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/scanners/typosquatting"
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/scanners/xssPrevention"
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/services/blockUrlAction"
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/services/visitAction"
)

const (
	blockURLActionKey = "block_url_action"
	visitActionKey    = "visit_action"
)

type Controller struct {
	scanners []Scanner
	actions  map[string]ActionGroup
	proxyStatus bool
	cacheCleared bool
	shutdown func(context.Context) error
	manageCache func(httplistener.CacheCommand)
}

type Scanner interface {
	Scan(r *http.Request) (blocked bool, reasons []string)
}

type ActionGroup interface {
	Name() string
}

type BlockActionGroup interface {
	ActionGroup
	BlockUrl(ctx context.Context, domain string, schedules []string) error
	GetAllBlockedURL(ctx context.Context) ([]string, error)
	GetBlockedURL(ctx context.Context, domain string) ([]string, error)
}

type VisitActionGroup interface {
	ActionGroup
	RegisterVisit(ctx context.Context, req *http.Request) error
}

func New(db *sql.DB) *Controller {
	scanners := []Scanner{
		xssPrevention.New(),
		blockURL.New(blockUrlDBService.NewBlockedDomainsDBService(db)),
		typosquatting.New(visitDBService.NewTypoSquattingDBService(db)),
		phishingPrevention.New(phishingDBService.NewPhishingDBService(db)),
	}
	actions := map[string]ActionGroup{
		blockURLActionKey: blockUrlAction.New(blockUrlDBService.NewBlockActionDomainsDBService(db)),
		visitActionKey:    visitAction.New(visitDBService.NewVisitActionDBService(db), phishingDBService.NewPhishingDBService(db)),
	}
	return &Controller{
		scanners: scanners,
		actions:  actions,
		proxyStatus: false,
		shutdown: nil,
		manageCache: nil,

	}
}

func (c* Controller) isProxyRunning() bool {
	return c.proxyStatus
}

func (c *Controller) isCacheCleared() bool {
	return c.cacheCleared
}

func (c *Controller) updateShutdownFunction(shutdown func(context.Context) error) {
	c.shutdown = shutdown
}

func (c *Controller) updateCacheFunction(manageCache func(httplistener.CacheCommand)) {
	c.manageCache = manageCache
}

func (c* Controller) updateProxyStatus(update bool) {
	c.proxyStatus = update
}

func (c* Controller) updateCacheCleared(update bool) {
	c.cacheCleared = update
}

func(c *Controller) runProxy() error {
	shutdown, manageCache, err := httplistener.ScanHTTPNetwork(c)
	if err != nil {
		return err
	}
	c.updateShutdownFunction(shutdown)
	c.updateCacheFunction(manageCache)
	c.updateProxyStatus(true)
	c.updateCacheCleared(true)
	fmt.Println("Server started on 127.0.0.1:4444")
	return nil
}

func (c* Controller) clearCache(targets []string) error{
	if c.manageCache == nil {
		return fmt.Errorf("proxy is offline, or did not start correctly")
	}
	if targets == nil{ //clear all
		c.manageCache(httplistener.CacheCommand{ClearAll: true})
	} else {
		c.manageCache(httplistener.CacheCommand{DeleteDomains: targets})
	}
	return nil
}

func (c * Controller) shutdownProxy() error {
	if c.shutdown == nil {
		return fmt.Errorf("proxy is offline, or did not start correctly")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	err := c.shutdown(ctx)
	cancel()
	if err != nil {
		return err
	}
	c.updateCacheFunction(nil)
	c.updateShutdownFunction(nil)
	c.updateProxyStatus(false)
	c.updateCacheCleared(true)
	fmt.Println("Server stopped.")
	return  nil
}

func (c *Controller) fetchBlockedDomains() ([]dto.BlockedDomainResponse, error) {
	blockGroup, err := c.blockActionGroup()
	if err != nil {
		return nil, fmt.Errorf("failed to get block action group: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	domains, err := blockGroup.GetAllBlockedURL(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get blocked URLs: %w", err)
	}

	output := make([]dto.BlockedDomainResponse, 0, len(domains))

	for _, domain := range domains {
		schedules, err := blockGroup.GetBlockedURL(ctx, domain)
		if err != nil {
			return nil, fmt.Errorf("failed to get schedules for url %q: %w", domain, err)
		}

		output = append(output, dto.BlockedDomainResponse{
			Domain:         domain,
			SchedulesCount: len(schedules),
			CreatedAt:      "",
			Schedules:      parseScheduleLinesToResponses(schedules),
		})
	}

	return output, nil
}

func (c *Controller) fetchBlockedDomain(domain string) (*dto.BlockedDomainResponse, error) {
	blockGroup, err := c.blockActionGroup()
	if err != nil {
		return nil, fmt.Errorf("failed to get block action group: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	schedules, err := blockGroup.GetBlockedURL(ctx, domain)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedules for url %q: %w", domain, err)
	}

	blockedDomain := &dto.BlockedDomainResponse{
		Domain:         domain,
		SchedulesCount: len(schedules),
		CreatedAt:      "",
		Schedules:      parseScheduleLinesToResponses(schedules),
	}

	return blockedDomain, nil
}

func (c *Controller) blockDomain(req dto.BlockedDomainRequest) (*dto.BlockedDomainResponse, error) {
	blockGroup, err := c.blockActionGroup()
	if err != nil {
		return nil, fmt.Errorf("failed to get block action group: %w", err)
	}

	domainSchedules, err := mapper.ToDomainSchedules(req)
	if err != nil {
		return nil, fmt.Errorf("invalid schedules: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = blockGroup.BlockURL(ctx, req.Domain, domainSchedules)
	if err != nil {
		return nil, fmt.Errorf("failed to block domain: %w", err)
	}

	if c.manageCache != nil {
		c.manageCache(httplistener.CacheCommand{
			DeleteDomains: []string{req.Domain},
		})
	}
	var schedules []dto.ScheduleResponse
	schedules = dto.RequestToScheduleResponses()
	return &dto.BlockedDomainResponse{
		Domain:         req.Domain,
		SchedulesCount: len(req.Schedules),
		CreatedAt:      req.CreatedAt,
		Schedules:      dto.ScheduleRequestToScheduleResponses(req.Schedules),
	}, nil
}

func (c *Controller) RunCLI() {
	reader := bufio.NewReader(os.Stdin)

	var shutdown func(context.Context) error
	var manageCache func(httplistener.CacheCommand)
	serverRunning := false

	for {
		fmt.Print(
			"<1> Passive Scan of Network\n" +
				"<2> Write block URL's\n" +
				"<3> Read blocked URL's\n" +
				"<4> Get History of Visited Domain\n" +
				"<5> Stop HTTP Server\n" +
				"<6> Clear session cache\n",
		)

		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("failed to read menu option: %v", err)
			return
		}

		switch strings.TrimSpace(line) {
		case "1":
			if serverRunning {
				fmt.Println("Server already running.")
				continue
			}

			shutdown, manageCache, err = httplistener.ScanHTTPNetwork(c)
			if err != nil {
				log.Printf("failed to start server: %v", err)
				continue
			}

			serverRunning = true
			fmt.Println("Server started on 127.0.0.1:4444")

		case "2":
			schedules, domain, err := readBlockInput(reader)
			if err != nil {
				log.Printf("failed to parse block input: %v", err)
				continue
			}

			blockGroup, err := c.blockActionGroup()
			if err != nil {
				log.Printf("failed to get block action group: %v", err)
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			err = blockGroup.BlockUrl(ctx, domain, schedules)
			cancel()
			if err != nil {
				log.Printf("failed to block domain %q: %v", domain, err)
				continue
			}

			if manageCache != nil {
				manageCache(httplistener.CacheCommand{DeleteDomains: []string{domain}})
			}

		case "3":
			fmt.Println("Enter a domain or skip to view all blocked domains:")

			line, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("failed to read domain input: %v", err)
				return
			}

			blockGroup, err := c.blockActionGroup()
			if err != nil {
				log.Printf("failed to get block action group: %v", err)
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

			if line != "\n" {
				domain := strings.TrimSpace(strings.ToLower(line))
				schedules, err := blockGroup.GetBlockedURL(ctx, domain)
				cancel()
				if err != nil {
					log.Printf("failed to get blocked URL %q: %v", domain, err)
					continue
				}
				fmt.Print(formatSchedules(schedules))
			} else {
				domains, err := blockGroup.GetAllBlockedURL(ctx)
				cancel()
				if err != nil {
					log.Printf("failed to get blocked URLs: %v", err)
					continue
				}
				fmt.Print(formatBlockedDomains(domains))
			}

		case "5":
			if !serverRunning {
				fmt.Println("Server not running.")
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			if err := shutdown(ctx); err != nil {
				log.Printf("failed to stop server: %v", err)
			}
			cancel()

			serverRunning = false
			fmt.Println("Server stopped.")

		case "6":
			if !serverRunning {
				fmt.Println("Server has to be running to clean its cache.")
				continue
			}

			manageCache(httplistener.CacheCommand{ClearAll: true})

		default:
			fmt.Println("Not implemented yet.")
		}
	}
}

func (c *Controller) InspectRequest(req *http.Request) (bool, string) {
	var reasons []string
	allowed := true

	for _, scanner := range c.scanners {
		blocked, scannerReasons := scanner.Scan(req)
		reasons = append(reasons, scannerReasons...)
		if blocked {
			allowed = false
		}
	}

	if allowed {
		if err := c.registerVisit(req); err != nil {
			log.Printf("failed to register visit: %v", err)
		}
	}

	log.Printf("controller scan result: allowed=%t reasons=%v", allowed, reasons)
	return allowed, formatScanReasons(reasons)
}

func (c *Controller) registerVisit(req *http.Request) error {
	group, ok := c.actions[visitActionKey]
	if !ok {
		return fmt.Errorf("visit action group not found")
	}

	visitGroup, ok := group.(VisitActionGroup)
	if !ok {
		return fmt.Errorf("invalid visit action group type")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := visitGroup.RegisterVisit(ctx, req); err != nil {
		return fmt.Errorf("register visit: %w", err)
	}

	return nil
}

func (c *Controller) blockActionGroup() (BlockActionGroup, error) {
	group, ok := c.actions[blockURLActionKey]
	if !ok {
		return nil, fmt.Errorf("block action group not found")
	}

	blockGroup, ok := group.(BlockActionGroup)
	if !ok {
		return nil, fmt.Errorf("invalid block action group type")
	}

	return blockGroup, nil
}

func formatScanReasons(reasons []string) string {
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

func readYesNo(reader *bufio.Reader) (bool, error) {
	for {
		response, err := reader.ReadString('\n')
		if err != nil {
			return false, fmt.Errorf("read yes/no response: %w", err)
		}

		response = strings.ToLower(strings.TrimSpace(response))
		switch response {
		case "y", "yes":
			return true, nil
		case "n", "no":
			return false, nil
		default:
			fmt.Print("Please enter [Yes/No]:")
		}
	}
}

func readBlockInput(reader *bufio.Reader) ([]string, string, error) {
	fmt.Print("Introduce a domain you which to block with this format: e.g. google.com\n")

	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, "", fmt.Errorf("read domain input: %w", err)
	}
	domain := strings.TrimSpace(strings.ToLower(line))

	fmt.Printf("Do you wish to set a block schedule for this domain, %s? [Yes/No]\n", domain)

	response, err := readYesNo(reader)
	if err != nil {
		return nil, "", fmt.Errorf("read schedule confirmation: %w", err)
	}

	if !response {
		return nil, domain, nil
	}

	schedules, err := readSchedule(reader)
	if err != nil {
		return nil, "", fmt.Errorf("read schedule input: %w", err)
	}

	return schedules, domain, nil
}

func readSchedule(reader *bufio.Reader) ([]string, error) {
	var schedules []string

	for {
		fmt.Print("Enter schedule as: <timestamp> <timestamp> <weekday>\nUse - to skip a field. Press Enter on an empty line to finish:\n")

		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("read schedule line: %w", err)
		}

		line = strings.TrimSpace(line)
		if line == "" {
			return schedules, nil
		}

		schedules = append(schedules, line)
	}
}

func formatSchedules(schedules []string) string {
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

func formatBlockedDomains(domains []string) string {
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

func parseSchedules(lines []string) []Schedule {
	output := make([]Schedule, 0, len(lines))

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Fields(line)

		schedule := Schedule{
			ID: fmt.Sprintf("%d", i+1),
		}

		if len(parts) > 0 && parts[0] != "-" {
			schedule.StartTime = parts[0]
		}

		if len(parts) > 1 && parts[1] != "-" {
			schedule.EndTime = parts[1]
		}

		if len(parts) > 2 && parts[2] != "-" {
			schedule.Weekday = parts[2]
		}

		output = append(output, schedule)
	}
	return output
}