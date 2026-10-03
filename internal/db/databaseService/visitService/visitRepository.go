package visitRepository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/db/dbmodel"
)

type VisitRepository struct{
	db *sql.DB
}

func NewVisitRepository(db *sql.DB) *VisitRepository{
	return &VisitRepository{
		db: db,
	}
}

func (v *VisitRepository) GetVisitedDomains(ctx context.Context) ([]dbmodel.Domain, error) {
	domains, err := fetchDomains(v.db, ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch domains from database: %s", err)
	}
	
	return domains, nil
}

func (v *VisitRepository) PushDomain(ctx context.Context, domain string) error {
	const q = `
		INSERT INTO visitedDomains (domain, time)
		VALUES (?, ?)
		ON DUPLICATE KEY UPDATE time = VALUES(time)
	`
	_, err := v.db.ExecContext(ctx, q, domain, time.Now())
	if err != nil {
		return fmt.Errorf("push domain: %w", err)
	}
	return nil
}

func fetchDomains(db *sql.DB, ctx context.Context) ([]dbmodel.Domain, error) {
	const q = `
		SELECT id, domain, time
		FROM visitedDomains
		ORDER BY time DESC
		LIMIT 500
	`

	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("fetch domains: %w", err)
	}
	defer rows.Close()

	var visited []dbmodel.Domain
	for rows.Next() {
		var d dbmodel.Domain
		if err := rows.Scan(&d.ID, &d.Domain, &d.Time); err != nil {
			return nil, fmt.Errorf("scan domain row: %w", err)
		}
		visited = append(visited, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}
	return visited, nil
}