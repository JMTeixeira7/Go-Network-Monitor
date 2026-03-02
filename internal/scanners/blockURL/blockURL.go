package blockURL

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/db/databaseService/blockUrlDBService"
)

type BlockedListDBService interface {
    IsDomainBlockedNow(ctx context.Context, domain string, now time.Time, day time.Weekday) (blocked bool, err error)
}

type Block struct{
	db_serivce BlockedListDBService
}

func New(b *blockUrlDBService.BlockUrlDBService) *Block {
	return &Block{
		db_serivce: b,
	}
}

func (b *Block) Scan(r *http.Request) (res bool, reasons []string) {
	ctx := r.Context()
	block, err := b.db_serivce.IsDomainBlockedNow(ctx, r.URL.Host, time.Now(), time.Now().Weekday())
	if err != nil {
		fmt.Printf("Error while using Blocked Domains DataBase service: %s\n", err)
		return true, nil
	}
	if !block{
		return block, reasons
	}
	return block, forgeScanMessage(r.URL.Host)
}

func forgeScanMessage(blockedDomain string) ([]string) {
	reasons := make([]string, 0)
	return  append(reasons, fmt.Sprintf("Request's domain, %s, is blocked\n", blockedDomain))
}