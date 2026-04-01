package blockUrlAction

import (
	"context"
	"fmt"

	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/model"
)

type dbservice interface {
	BlockUrlDB(ctx context.Context, domain string, schedules []*model.Schedule) error
	GetAllBlockedURL(ctx context.Context) ([]string, error)
	GetBlockedURL(ctx context.Context, domain string) ([]*model.Schedule, error)
}

type Service struct {
	dbservice dbservice
}

func New(dbservice dbservice) *Service {
	return &Service{dbservice: dbservice}
}

func (s *Service) Name() string {
	return "block_url_action"
}

func (s *Service) BlockUrl(ctx context.Context, domain string, schedules []*model.Schedule) error {
	if err := s.dbservice.BlockUrlDB(ctx, domain, schedules); err != nil {
		return fmt.Errorf("store blocked URL: %w", err)
	}

	return nil
}

func (s *Service) GetAllBlockedURL(ctx context.Context) ([]string, error) {
	domains, err := s.dbservice.GetAllBlockedURL(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all blocked URLs: %w", err)
	}

	return domains, nil
}

func (s *Service) GetBlockedURL(ctx context.Context, domain string) ([]*model.Schedule, error) {
	schedules, err := s.dbservice.GetBlockedURL(ctx, domain)
	if err != nil {
		return nil, fmt.Errorf("get blocked URL %q: %w", domain, err)
	}
	return schedules, nil
}
