package visitAction

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/db/databaseService/phishingDBService"
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/db/databaseService/visitDBService"
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/model"
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/resources/credentialsParser"
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/security"
)

type visitUrlDBService interface {
	PushDomain(ctx context.Context, domain string) error
}

type credsDBService interface {
	PushCredentials(ctx context.Context, cred *model.Credentials, domain string) error
}

type VisitService struct {
	db_visit_service  visitUrlDBService
	db_creds_service credsDBService
	fp_service *security.Fingerprinter

}

func New(db_visit_service *visitDBService.VisitActionDBService, db_creds_service *phishingDBService.PhishingDBService ) *VisitService {
	fp_service := security.NewFingerprinter([]byte(os.Getenv("FINGERPRINT_SEED")))
	return &VisitService{
		db_visit_service: db_visit_service,
		fp_service: fp_service,
		db_creds_service: db_creds_service,

	}
}

func (b *VisitService) Name() string {
	return "visit_service"
}

func (v *VisitService) RegisterVisit(ctx context.Context, req *http.Request) error {
	err := v.db_visit_service.PushDomain(ctx, req.URL.Host)
	if err != nil {
		return fmt.Errorf("Error while pushing current request to database: %w\n", err)
	}
	res, err := v.pushCredentials(ctx, req)
	if  err != nil {
		return fmt.Errorf("Error while using credentials DataBase service: %w\n", err)
	}
	if res {
		fmt.Println("Credentials pushed into db successfully")
	}
	return nil
}

func (v *VisitService) pushCredentials(ctx context.Context, req *http.Request) (bool, error) {
	if req.Method != http.MethodPost {
		return false, nil
	}
	creds := inspectRequest(req, v.fp_service)
	if creds==nil {
		return false, nil
	}
	err := v.db_creds_service.PushCredentials(ctx, creds, req.URL.Host)
	if err != nil {
		return false, err
	}
	return true, nil
}

func inspectRequest(req *http.Request, fp_service *security.Fingerprinter) *model.Credentials {
	err := req.ParseForm()
	if err != nil {
		fmt.Println("Error while parsing the req Form")
		return nil
	}
	email, username, password := credentialsParser.ExtractCredentialFields(req.Form)

	fmt.Printf("username: %s, email: %s, password: %s\n", username, email, password)
	if !(username == "" && email == "") && password != "" {
		return model.CreateCredentials(email, username, password, fp_service)
	}
	return nil
}
