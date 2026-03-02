package blockUrlDBService

import (
	"context"
	"database/sql"
	"time"
)

type BlockUrlDBService struct {
	db *sql.DB
} 

func NewBlockedDomainsDBService(db *sql.DB) *BlockUrlDBService {
	return &BlockUrlDBService{
		db: db,
	}
}

func (b *BlockUrlDBService) IsDomainBlockedNow(ctx context.Context, domain string, now time.Time, day time.Weekday) (blocked bool, err error) {
	panic("not yet implemented")
}


