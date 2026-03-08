package phishingPrevention

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/db/databaseService/phishingDBService"
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/model"
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/security"
)

type DBService interface {
	CheckForPhishing(ctx context.Context, cred *model.Credentials, domain string) (bool, *string, error)
	PushCredentials(ctx context.Context, cred *model.Credentials, domain string) error
}

type PhishingPrev struct {
	db_service DBService
	fp_service *security.Fingerprinter
}

func New(p *phishingDBService.PhishingDBService) *PhishingPrev{
	fp_service := security.NewFingerprinter([]byte(os.Getenv("FINGERPRINT_SEED")))
	return &PhishingPrev{
		db_service: p,
		fp_service: fp_service,
	}
}

func (p *PhishingPrev) Scan(req *http.Request) (bool, []string) {
	if req.Method != "POST" {
		return false, nil
	}
	creds := inspectRequest(req, p.fp_service)
	if creds==nil {
		return false, nil
	}
	ctx := req.Context()
	phishing, reason, err := p.db_service.CheckForPhishing(ctx, creds, req.URL.Host)
	if err != nil {
		fmt.Printf("Error while using phishing DataBase service: %s\n", err)
		return false, nil
	}
	err = p.db_service.PushCredentials(ctx, creds, req.URL.Host)
	if err != nil {
		fmt.Printf("Error while using phishing DataBase service: %s\n", err)
	}
	reasons := []string{}
	reasons = append(reasons, *reason)
	return phishing, reasons
}

func inspectRequest(req *http.Request, fp_service *security.Fingerprinter) *model.Credentials {
	err := req.ParseForm()
	if err != nil {
		return nil
	}
	username := req.FormValue("username")
	email := req.FormValue("email")
	password := req.FormValue("password")
	if !(username == "" && email == "") || password != "" {
		return model.CreateCredentials(email, username, password, fp_service)
	}
	return nil	//No password or no id, no leak
}