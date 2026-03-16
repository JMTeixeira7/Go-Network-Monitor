package phishingDBService

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/model"
)

type PhishingDBService struct {
	db *sql.DB
}

func NewPhishingDBService(db *sql.DB) *PhishingDBService {
	return &PhishingDBService{
		db: db,
	}
}

func (p *PhishingDBService) CheckForPhishing(ctx context.Context, cred *model.Credentials, domain string) (bool, *string, error) {
	phishing, err := fetchPhishingInstance(ctx, cred, domain, p.db)
	if err != nil {
		return false, nil, fmt.Errorf("Failed to fetch credentials instance from database: %s", err)
	}
	if phishing != nil {
		return true, phishing, nil
	}
	return false, nil, nil
}

func (p *PhishingDBService) PushCredentials(ctx context.Context, cred *model.Credentials, domain string) error {
	const q = `
		INSERT INTO credentials (domain_key, username, fingerprint)
		VALUES (?, ?, ?)
	`
	
	domain_key, err := fetchDomainIDInstance(ctx, domain, p.db)
	if err != nil {
		return fmt.Errorf("Failed to push credentials into table:\n %s", err)
	}
	_, err = p.db.ExecContext(ctx, q, domain_key, cred.Username, cred.Fingerprint)
	if err != nil {
		return fmt.Errorf("push domain: %w", err)
	}
	return nil
}

func fetchPhishingInstance(ctx context.Context, cred *model.Credentials, domain string, db *sql.DB) (*string, error) {
	const q = `
		SELECT v.domain
		FROM credentials c
		JOIN visitedDomains v ON v.id = c.domain_key
		WHERE c.username = ? AND c.fingerprint = ? AND v.domain <> ?
		ORDER BY v.time DESC
		LIMIT 1;
	`
	fmt.Printf("Phishing test attempt-> credntials: %+v", cred)
	var legitDomain string
	err := db.QueryRowContext(ctx, q, cred.Username, cred.Fingerprint, domain).Scan(&legitDomain)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &legitDomain, nil
}

func fetchDomainIDInstance(ctx context.Context, domain string, db *sql.DB) (*int, error) {

	const q = `
		SELECT id
		FROM visitedDomains
		WHERE domain = ?
		LIMIT 1
	`
	var domain_key int
	err := db.QueryRowContext(ctx, q, domain).Scan(&domain_key)
	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("Could not fetch a domain instance for the given credentials (Unwanted behaviour): %s", domain)
	}
	if err != nil {
		return  nil, err
	}
	return &domain_key, nil
		
}