package blockUrlAction

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/db/databaseService/blockUrlDBService"
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/model"
)

type BlockActionUrlDBService interface {
	BlockUrlDB(ctx context.Context, domain string, schedules []*model.Schedule) error
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

func (b *BlockURLService) BlockUrl(ctx context.Context, domain string, raw_schedules []string) error {
	parsed_schedules, err := parseSchedules(raw_schedules)
	if err !=nil {
		return fmt.Errorf("Error while parsing the Schedules, %w", err)
	}
	err = b.db_service.BlockUrlDB(ctx, domain, parsed_schedules)
	if err != nil {
		return err
	}
	return nil
}

func (b *BlockURLService) GetAllBlockedURL(ctx context.Context) ([]string, error) {
	panic("not yet implemented")
}

func (b *BlockURLService) GetBlockedURL(ctx context.Context, domain string) ([]string, error) {
	panic("not yet implemented")
}


func parseSchedules(lines []string) ([]*model.Schedule, error) {
	var schedules []*model.Schedule

	for i, line := range lines {
		schedule, err := parseScheduleLine(line)
		if err != nil {
			return nil, fmt.Errorf("error parsing schedule at line %d: %w", i+1, err)
		}

		if schedule != nil {
			schedules = append(schedules, schedule)
		}
	}
	return schedules, nil
}

func decodeSchedules(parsed_schedules []model.Schedule) (decoded_schedule []string) {
	panic("not implemented")
}

func parseScheduleLine(line string) (*model.Schedule, error) {
	fields := strings.Fields(line)

	if len(fields) != 3 {
		return nil, fmt.Errorf("schedule must have exactly 3 fields: <timestamp> <timestamp> <weekday>")
	}

	startTime, err := parseTimestamp(fields[0])
	if err != nil {
		return nil, err
	}

	endTime, err := parseTimestamp(fields[1])
	if err != nil {
		return nil, err
	}

	weekday, err := parseWeekday(fields[2])
	if err != nil {
		return nil, err
	}

	schedule, err := model.CreateSchedule(startTime, endTime, weekday)
	if err != nil {
		return nil, err
	}

	return schedule, nil
}

func parseWeekday(s string) (*time.Weekday, error) {
	if s == "-" || s == "" {
		return nil, nil
	}

	switch strings.ToLower(s) {
	case "sunday":
		w := time.Sunday
		return &w, nil
	case "monday":
		w := time.Monday
		return &w, nil
	case "tuesday":
		w := time.Tuesday
		return &w, nil
	case "wednesday":
		w := time.Wednesday
		return &w, nil
	case "thursday":
		w := time.Thursday
		return &w, nil
	case "friday":
		w := time.Friday
		return &w, nil
	case "saturday":
		w := time.Saturday
		return &w, nil
	default:
		return nil, fmt.Errorf("invalid weekday: %s", s)
	}
}

func parseTimestamp(s string) (*time.Time, error) {
	if s == "-" || s == "" {
		return nil, nil
	}

	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp %q: %w", s, err)
	}
	return &t, nil
}
