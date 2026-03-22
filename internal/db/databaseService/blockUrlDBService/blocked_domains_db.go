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
	return &BlockActionUrlDBService{db: db}
}

func NewBlockedDomainsDBService(db *sql.DB) *BlockUrlDBService {
	return &BlockUrlDBService{db: db}
}

func (a *BlockActionUrlDBService) BlockUrlDB(ctx context.Context, domain string, schedules []*model.Schedule) error {
	dbSchedules, err := toDBSchedules(schedules)
	if err != nil {
		return fmt.Errorf("convert schedules for domain %q: %w", domain, err)
	}

	if err := blockURLTransaction(a.db, ctx, domain, dbSchedules); err != nil {
		return fmt.Errorf("store blocked domain %q: %w", domain, err)
	}

	return nil
}

func (a *BlockActionUrlDBService) GetAllBlockedURL(ctx context.Context) ([]string, error) {
	domains, err := fetchBlockedDomains(a.db, ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch blocked domains: %w", err)
	}
	return domains, nil
}

func (a *BlockActionUrlDBService) GetBlockedURL(ctx context.Context, domain string) ([]*model.Schedule, error) {
	dbSchedules, err := fetchBlockedDomainSchedules(a.db, ctx, domain)
	if err != nil {
		return nil, fmt.Errorf("fetch schedules for domain %q: %w", domain, err)
	}

	modelSchedules, err := toModelSchedules(dbSchedules)
	if err != nil {
		return nil, fmt.Errorf("convert schedules for domain %q: %w", domain, err)
	}

	return modelSchedules, nil
}

func (b *BlockUrlDBService) IsDomainBlockedNow(ctx context.Context, domain string, now *time.Time, day *time.Weekday) (bool, error) {
	dbSchedules, err := fetchBlockedDomainSchedules(b.db, ctx, domain)
	if err != nil {
		return false, fmt.Errorf("fetch schedules for domain %q: %w", domain, err)
	}
	if len(dbSchedules) == 0 {
		return false, nil
	}

	return dbmodel.AnyScheduleMatches(dbSchedules, now, day), nil
}

func fetchBlockedDomains(db *sql.DB, ctx context.Context) ([]string, error) {
	const q = `
		SELECT domain
		FROM blockedDomains
	`

	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("query blocked domains: %w", err)
	}
	defer rows.Close()

	var domains []string
	for rows.Next() {
		var domain string
		if err := rows.Scan(&domain); err != nil {
			return nil, fmt.Errorf("scan blocked domain: %w", err)
		}
		domains = append(domains, domain)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate blocked domains rows: %w", err)
	}

	return domains, nil
}

func blockURLTransaction(db *sql.DB, ctx context.Context, domain string, schedules []*dbmodel.Schedule) error {
	const q1 = `
		INSERT INTO blockedDomains
		(domain)
		VALUES (?)
		ON DUPLICATE KEY UPDATE id = LAST_INSERT_ID(id)
	`
	const q2 = `
		INSERT INTO schedule
		(blocked_domain_key, start_time, end_time, weekday, timezone)
		VALUES
		(?, ?, ?, ?, ?)
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
		return fmt.Errorf("insert blocked domain: %w", err)
	}

	blockedDomainKey, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("get blocked domain id: %w", err)
	}

	for _, schedule := range schedules {
		start, end, weekday, timezone := schedule.SQLValues()

		_, err := tx.ExecContext(ctx, q2, blockedDomainKey, start, end, weekday, timezone)
		if err != nil {
			return fmt.Errorf("insert schedule: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	committed = true
	return nil
}

func fetchBlockedDomainSchedules(db *sql.DB, ctx context.Context, domain string) ([]*dbmodel.Schedule, error) {
	const q = `
		SELECT s.start_time, s.end_time, s.weekday, s.timezone
		FROM blockedDomains b
		JOIN schedule s ON b.id = s.blocked_domain_key
		WHERE b.domain = ?
	`

	rows, err := db.QueryContext(ctx, q, domain)
	if err != nil {
		return nil, fmt.Errorf("query schedules for domain %q: %w", domain, err)
	}
	defer rows.Close()

	var schedules []*dbmodel.Schedule
	for rows.Next() {
		schedule, err := scanScheduleRow(rows)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate schedules rows for domain %q: %w", domain, err)
	}

	return schedules, nil
}

func scanScheduleRow(rows *sql.Rows) (*dbmodel.Schedule, error) {
	var startStr, endStr sql.NullString
	var weekdayValue, timezoneValue sql.NullInt64

	if err := rows.Scan(&startStr, &endStr, &weekdayValue, &timezoneValue); err != nil {
		return nil, fmt.Errorf("scan schedule row: %w", err)
	}

	var startTime *dbmodel.Clock
	if startStr.Valid {
		clock, err := dbmodel.ParseClockString(startStr.String)
		if err != nil {
			return nil, fmt.Errorf("parse start time %q: %w", startStr.String, err)
		}
		startTime = clock
	}

	var endTime *dbmodel.Clock
	if endStr.Valid {
		clock, err := dbmodel.ParseClockString(endStr.String)
		if err != nil {
			return nil, fmt.Errorf("parse end time %q: %w", endStr.String, err)
		}
		endTime = clock
	}

	var weekday *time.Weekday
	if weekdayValue.Valid {
		w := time.Weekday(weekdayValue.Int64)
		weekday = &w
	}

	var timezone *int
	if timezoneValue.Valid {
		tz := int(timezoneValue.Int64)
		timezone = &tz
	}

	schedule, err := dbmodel.NewSchedule(startTime, endTime, weekday, timezone)
	if err != nil {
		return nil, fmt.Errorf("build db schedule from row: %w", err)
	}

	return schedule, nil
}

func toModelSchedules(dbSchedules []*dbmodel.Schedule) ([]*model.Schedule, error) {
	modelSchedules := make([]*model.Schedule, 0, len(dbSchedules))

	for i, schedule := range dbSchedules {
		modelSchedule, err := toModelSchedule(schedule)
		if err != nil {
			return nil, fmt.Errorf("convert schedule %d to model: %w", i+1, err)
		}
		modelSchedules = append(modelSchedules, modelSchedule)
	}

	return modelSchedules, nil
}

func toModelSchedule(schedule *dbmodel.Schedule) (*model.Schedule, error) {
	if schedule == nil {
		return nil, nil
	}

	return model.CreateScheduleFromDB(
		toModelClock(schedule.StartTime()),
		toModelClock(schedule.EndTime()),
		schedule.Weekday(),
		schedule.Timezone(),
	)
}

func toModelClock(clock *dbmodel.Clock) *model.Clock {
	if clock == nil {
		return nil
	}

	return model.CreateClock(clock.Hour(), clock.Minute(), clock.Second())
}

func toDBSchedules(modelSchedules []*model.Schedule) ([]*dbmodel.Schedule, error) {
	if len(modelSchedules) == 0 {
		return []*dbmodel.Schedule{nil}, nil
	}

	dbSchedules := make([]*dbmodel.Schedule, 0, len(modelSchedules))
	for i, schedule := range modelSchedules {
		dbSchedule, err := toDBSchedule(schedule)
		if err != nil {
			return nil, fmt.Errorf("convert schedule %d to db model: %w", i+1, err)
		}
		dbSchedules = append(dbSchedules, dbSchedule)
	}

	return dbSchedules, nil
}

func toDBSchedule(schedule *model.Schedule) (*dbmodel.Schedule, error) {
	if schedule == nil {
		return nil, nil
	}

	startTime, err := toDBClock(schedule.StartTime())
	if err != nil {
		return nil, fmt.Errorf("convert start time: %w", err)
	}

	endTime, err := toDBClock(schedule.EndTime())
	if err != nil {
		return nil, fmt.Errorf("convert end time: %w", err)
	}

	dbSchedule, err := dbmodel.NewSchedule(
		startTime,
		endTime,
		schedule.Weekday(),
		schedule.Timezone(),
	)
	if err != nil {
		return nil, fmt.Errorf("validate schedule for persistence: %w", err)
	}

	return dbSchedule, nil
}

func toDBClock(clock *model.Clock) (*dbmodel.Clock, error) {
	if clock == nil {
		return nil, nil
	}
	return dbmodel.NewClock(clock.Hour(), clock.Minute(), clock.Second())
}
