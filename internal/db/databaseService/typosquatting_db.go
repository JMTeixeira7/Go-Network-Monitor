package databaseservice

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/db/model"
)

type TyposquattingDBService struct{
	db *sql.DB
}

func New(db *sql.DB) *TyposquattingDBService{
	return &TyposquattingDBService{
		db: db,
	}
}

func (t *TyposquattingDBService) GetVisitedDomains(ctx context.Context) ([]string, error) {
	domains, err := fetchDomains(t.db, ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch domains from database: %s", err)
	}
	str_domains := []string{}
	for i := 0; i < len(domains); i++ {
		str_domains = append(str_domains, domains[i].Domain)
	}
	return str_domains, nil
}

func (t *TyposquattingDBService) PushDomain(ctx context.Context, domain string) error {
	const q = `
		INSERT INTO visitedDomains (domain, time)
		VALUES (?, ?)
	`
	_, err := t.db.ExecContext(ctx, q, domain, time.Now())
	if err != nil {
		return fmt.Errorf("push domain: %w", err)
	}
	return nil
}

func fetchDomains(db *sql.DB, ctx context.Context) ([]model.Domain, error) {
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

	var visited []model.Domain
	for rows.Next() {
		var d model.Domain
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