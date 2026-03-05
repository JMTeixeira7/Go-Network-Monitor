package blockUrlDBService

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/db/model"
)

type BlockUrlDBService struct {
	db *sql.DB
}

func NewBlockedDomainsDBService(db *sql.DB) *BlockUrlDBService {
	return &BlockUrlDBService{
		db: db,
	}
}

func (b *BlockUrlDBService) IsDomainBlockedNow(ctx context.Context, domain string, now *time.Time, day *time.Weekday) (blocked bool, err error) {
	schedules, err := fetchBlockedDomainSchedules(b.db, ctx, domain)
	if err != nil {
		return false, fmt.Errorf("Error while checking if domain %s is blocked: %w", domain, err)
	}
	if schedules != nil {
		blocked = isCurrentlyBlocked(schedules, now, day)
		return blocked, nil
	}
	return true, nil
}

func fetchBlockedDomainSchedules(db *sql.DB, ctx context.Context, domain string) (schedules []model.Schedule, err error) {
	const q = `
		SELECT start_time, end_time, weekday
		FROM blockedDomains b
		JOIN schedule s
		WHERE b.domain = ? AND b.id = s.blocked_domain_key
	`
	rows, err := db.QueryContext(ctx, q, domain)
	if err != nil {
		return nil, fmt.Errorf("Error while fetching Blocked domain key: %w", err)
	}
	defer rows.Close()

	var schedule_rows []model.Schedule
	for rows.Next() {
		var s model.Schedule
		if err := rows.Scan(&s.Start_time, &s.End_time, &s.Weekday); err != nil {
			return nil, fmt.Errorf("Fail to scan schedule row: %w", err)
		}
		schedule_rows = append(schedule_rows, s)
	}
	err = rows.Err()
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil	//domain is blocked without explicit schedule
		}
		return nil, fmt.Errorf("Fail iterating rows: %w", err)
	}
	return schedule_rows, nil
}

func isCurrentlyBlocked(schedules []model.Schedule, now *time.Time, day *time.Weekday) bool {
	for i := range schedules {
		if schedules[i].Weekday != nil {
			if day != nil && day == schedules[i].Weekday && timeslotsIntersect(now, schedules[i].Start_time, schedules[i].End_time) {
				return true
			}
			continue
		}
		if timeslotsIntersect(now, schedules[i].Start_time, schedules[i].End_time) {
			return true
		}
	}
	return false
}

func timeslotsIntersect(now *time.Time, min *time.Time, max *time.Time) bool {
	if now == nil || min == nil || max == nil {
		return true
	}
	if max.Before(*min) || max.Equal(*min) {
		return false
	}
	if min.Before(*now) && max.After(*now) {
		return true
	}
	return false
}


