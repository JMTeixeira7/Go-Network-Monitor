package blockUrlAction

import (
	"context"
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/db/databaseService/blockUrlDBService"
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/model"
)

type BlockActionUrlDBService interface {
	BlockUrlDB(ctx context.Context, domain string, schedules []model.Schedule) error
	GetAllBlockedURL(ctx context.Context) ([]string, error)
	GetBlockedURL(ctx context.Context, domain string) ([]model.Schedule, error)

}
type BlockURLService struct {
	db_service  BlockActionUrlDBService
}

func New(db_service *blockUrlDBService.BlockActionUrlDBService) *BlockURLService {
	return &BlockURLService{
		db_service: db_service,
	}
}

func (b *BlockURLService) Name() string {
	return "block_url_action"
}

func (b *BlockURLService) BlockUrl(domain string, schedules []model.Schedule) error {
	panic("not yet implemented")
}

func (b *BlockURLService) GetAllBlockedURL() ([]string, error) {
	panic("not yet implemented")
}

func (b *BlockURLService) GetBlockedURL(domain string) ([]model.Schedule, error) {
	panic("not yet implemented")
}
