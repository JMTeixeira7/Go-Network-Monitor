package blockURL

import (
	"net/http"
	"context"
	"fmt"
	"time"

)

type BlockedListDBService interface {
    IsDomainBlockedNow(ctx context.Context, domain string, now time.Time, day time.Weekday) (blocked bool, err error)
}

type Block struct{
	db_serivce BlockedListDBService
}

func New(db_service BlockedListDBService) *Block {
	return &Block{
		db_serivce: db_service,
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