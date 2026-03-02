package phishingDBService

import (
	"context"
	"database/sql"

	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/model"
)

type PhishingDBService struct {
	db *sql.DB
}

func NewPhishingDBService(db *sql.DB) *PhishingDBService {
	return &PhishingDBService{
		db: db,
	}
}

func (p *PhishingDBService) CheckForPhishing(ctx context.Context, cred *model.Credentials, domain string) (bool, string, error) {
	panic("not yet implemented")
}

func (p *PhishingDBService) PushCredentials(ctx context.Context, cred *model.Credentials, domain string) error {
	panic("not yet implemented")
}