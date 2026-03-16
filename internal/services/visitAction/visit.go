package visitAction

import (
	"context"
	"fmt"
	"net/http"

	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/db/databaseService/visitDBService"
)

type visitUrlDBService interface {
	PushDomain(ctx context.Context, domain string) error
}

type VisitService struct {
	db_service  visitUrlDBService
}

func New(db_service *visitDBService.VisitActionDBService) *VisitService {
	return &VisitService{
		db_service: db_service,
	}
}

func (b *VisitService) Name() string {
	return "visit_service"
}

func (v *VisitService) RegisterVisit(ctx context.Context, req *http.Request) error {
	err := v.db_service.PushDomain(ctx, req.URL.Hostname())
	if err != nil {
		return fmt.Errorf("Error while pushing current request to database: %w\n", err)
	}
	return nil
}