package blockUrlDBService

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/db/dbmodel"
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/model"
)

type BlockUrlDBService struct {
	db *sql.DB
}

type BlockActionUrlDBService struct {
	db *sql.DB
}

func NewBlockActionDomainsDBService(db *sql.DB) *BlockActionUrlDBService {
	return &BlockActionUrlDBService{
		db: db,
	}
}

func NewBlockedDomainsDBService(db *sql.DB) *BlockUrlDBService {
	return &BlockUrlDBService{
		db: db,
	}
}


func (a *BlockActionUrlDBService) BlockUrlDB(ctx context.Context, domain string, schedules []*model.Schedule) error {
	var db_schedules []dbmodel.Schedule
	for _, s := range schedules {
		db_s := toDBSchedule(s)
		db_schedules = append(db_schedules, *db_s)
	}
	err := blockUrlTransaction(a.db, ctx, domain, db_schedules)
	if err != nil {
		return err
	}
	return nil
}

func (a *BlockActionUrlDBService) GetAllBlockedURL(ctx context.Context) ([]string, error) {
	blocked_domains, err := fetchBlockedDomains(a.db, ctx)
	if err != nil {
		return nil, fmt.Errorf("Error while fetching Blocked domains from db: %w", err)
	}
	return blocked_domains, nil
}

func (a *BlockActionUrlDBService) GetBlockedURL(ctx context.Context, domain string) ([]*model.Schedule, error) {
	schedules, err := fetchBlockedDomainSchedules(a.db, ctx, domain)
	if err != nil {
		return nil, fmt.Errorf("Error while fetching domain %s block schedules: %w", domain, err)
	}
	var model_schedules []*model.Schedule
	for _, s := range schedules {
		model_s, err := toModelSchedule(&s)
		if err != nil {
			return nil, fmt.Errorf("Error while parsing db Schedule, %w", err)
		}
		model_schedules = append(model_schedules, model_s)
	}
	return model_schedules, nil
}

func (b *BlockUrlDBService) IsDomainBlockedNow(ctx context.Context, domain string, now *time.Time, day *time.Weekday) (blocked bool, err error) {
	schedules, err := fetchBlockedDomainSchedules(b.db, ctx, domain)
	if err != nil {
		return false, fmt.Errorf("Error while fetching domain %s block schedules: %w", domain, err)
	}
	if schedules != nil {
		blocked = isCurrentlyBlocked(schedules, now, day)
		return blocked, nil
	}
	return true, nil
}

func fetchBlockedDomains(db *sql.DB, ctx context.Context) ([]string, error) {
	const q = `
		SELECT domain
		FROM blockedDomains
	`
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var blocked_domains []string
	for rows.Next() {
		var d string
		if err := rows.Scan(&d); err != nil {
			return nil, fmt.Errorf("Fail to scan schedule row: %w", err)
		}
		blocked_domains = append(blocked_domains, d)
	}
	err = rows.Err()
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil	// no blocked domains
		}
		return nil, fmt.Errorf("Fail iterating rows: %w", err)
	}
	return blocked_domains, nil
}

func blockUrlTransaction(db *sql.DB, ctx context.Context, domain string, schedules []dbmodel.Schedule) error {
	const q1 = `
		INSERT INTO blockedDomains
		(domain)
		VALUES
		(?)
	`
	const q2 = `
		INSERT INTO schedule
		(blocked_domain_key, start_time, end_time, weekday)
		VALUES
		(?, ?, ?, ?)
	`

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	res, err := tx.ExecContext(ctx, q1, domain)
	if err != nil {
		return fmt.Errorf("push block_domain: %w", err)
	}
	blocked_domain_key, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("get blocked domain id: %w", err)
	}
	for _, s := range schedules{
		var weekday_value any
		if s.Weekday == nil { weekday_value = nil
		} else { weekday_value = int(*s.Weekday)}

		_, err := tx.ExecContext(ctx, q2, blocked_domain_key, s.Start_time, s.End_time, weekday_value)
		if err != nil {
			return fmt.Errorf("push schedule: %w", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit transaction aborted: %w", err)
	}
	committed = true
	return nil
}
func fetchBlockedDomainSchedules(db *sql.DB, ctx context.Context, domain string) (schedules []dbmodel.Schedule, err error) {
	const q = `
		SELECT s.start_time, s.end_time, s.weekday
		FROM blockedDomains b
		JOIN schedule s ON b.id = s.blocked_domain_key
		WHERE b.domain = ?
	`

	rows, err := db.QueryContext(ctx, q, domain)
	if err != nil {
		return nil, fmt.Errorf("Error while fetching Blocked domain key: %w", err)
	}
	defer rows.Close()

	var schedule_rows []dbmodel.Schedule
	for rows.Next() {
		var s dbmodel.Schedule
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

func isCurrentlyBlocked(schedules []dbmodel.Schedule, now *time.Time, day *time.Weekday) bool {
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

func toModelSchedule(schedule *dbmodel.Schedule) (*model.Schedule, error) {
	if schedule == nil {
		return nil, nil
	}

	return model.CreateSchedule(
		schedule.Start_time,
		schedule.End_time,
		schedule.Weekday,
	)
}

func toDBSchedule(schedule *model.Schedule) *dbmodel.Schedule {
	if schedule == nil {
		return nil
	}

	return &dbmodel.Schedule{
		Start_time: schedule.StartTime(),
		End_time:   schedule.EndTime(),
		Weekday:    schedule.Weekday(),
	}
}
