package visitAction

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/model"
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/resources/credentialsParser"
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/security"
)

type visitStore interface {
	PushDomain(ctx context.Context, domain string) error
}

type credentialsStore interface {
	PushCredentials(ctx context.Context, cred *model.Credentials, domain string) error
}

type Service struct {
	visitStore       visitStore
	credentialsStore credentialsStore
	fingerprinter    *security.Fingerprinter
}

func New(visitStore visitStore, credentialsStore credentialsStore) *Service {
	fingerprinter := security.NewFingerprinter([]byte(os.Getenv("FINGERPRINT_SEED")))

	return &Service{
		visitStore:       visitStore,
		credentialsStore: credentialsStore,
		fingerprinter:    fingerprinter,
	}
}

func (s *Service) Name() string {
	return "visit_service"
}

func (s *Service) RegisterVisit(ctx context.Context, req *http.Request) error {
	if err := s.visitStore.PushDomain(ctx, req.URL.Host); err != nil {
		return fmt.Errorf("push visited domain: %w", err)
	}

	stored, err := s.pushCredentials(ctx, req)
	if err != nil {
		return fmt.Errorf("store credentials: %w", err)
	}

	if stored {
		log.Printf("stored credentials for host=%s", req.URL.Host)
	}

	return nil
}

func (s *Service) pushCredentials(ctx context.Context, req *http.Request) (bool, error) {
	if req.Method != http.MethodPost {
		return false, nil
	}

	creds, err := s.extractCredentials(req)
	if err != nil {
		return false, fmt.Errorf("extract credentials from request: %w", err)
	}
	if creds == nil {
		return false, nil
	}

	if err := s.credentialsStore.PushCredentials(ctx, creds, req.URL.Host); err != nil {
		return false, fmt.Errorf("push credentials: %w", err)
	}

	return true, nil
}

func (s *Service) extractCredentials(req *http.Request) (*model.Credentials, error) {
	if err := req.ParseForm(); err != nil {
		return nil, fmt.Errorf("parse request form: %w", err)
	}

	email, username, password := credentialsParser.ExtractCredentialFields(req.Form)

	if (username == "" && email == "") || password == "" {
		return nil, nil
	}

	if username != "" || email != "" {
		log.Printf("credential candidate detected for host=%s username=%q email=%q", req.URL.Host, username, email)
	}

	return model.CreateCredentials(email, username, password, s.fingerprinter), nil
}
