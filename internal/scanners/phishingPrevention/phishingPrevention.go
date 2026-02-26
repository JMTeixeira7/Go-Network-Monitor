package phishingPrevention

import(
	"fmt"
	"context"
	"http"

)

type phishingDBService interface {
	//TODO
}

type PhishingPrev struct {
	db_service phishingDBService
}

func New(db_service phishingDBService) *PhishingPrev{
	return &PhishingPrev{
		db_service: db_service,
	}
}

func (p *PhishingPrev) Scan(r *http.Request) (res bool, reasons []string) {
	panic("not yet implemented")
}