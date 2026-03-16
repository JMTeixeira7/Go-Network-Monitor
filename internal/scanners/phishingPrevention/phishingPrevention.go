package phishingPrevention

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"net/url"
	"strings"
	"unicode"

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
	if req.Method != http.MethodPost {
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
	reasons := forgeScanMessage(req.URL.Host, *reason)
	return phishing, reasons
}

func inspectRequest(req *http.Request, fp_service *security.Fingerprinter) *model.Credentials {
	err := req.ParseForm()
	if err != nil {
		fmt.Println("Error while parsing the req Form")
		return nil
	}
	email, username, password := extractCredentialFields(req.Form)

	fmt.Printf("username: %s, email: %s, password: %s\n", username, email, password)
	if !(username == "" && email == "") && password != "" {
		fmt.Println("creates credentials")
		return model.CreateCredentials(email, username, password, fp_service)
	}
	return nil	//No password or no id, no leak
}

func forgeScanMessage(blockedDomain string, legitDomain string) ([]string) {
	reasons := make([]string, 0)
	return  append(reasons, fmt.Sprintf(`Authentication request for domain, %s, was blocked due to
					 Phishing detected on Legit domain: %s\n`, blockedDomain, legitDomain))
}

func extractCredentialFields(form url.Values) (email, username, password string) {
	for rawKey, values := range form {
		if len(values) == 0 {
			continue
		}

		value := strings.TrimSpace(values[0])
		if value == "" {
			continue
		}

		key := normalizeFieldKey(rawKey)

		switch {
		case password == "" && looksLikePasswordKey(key):
			password = value

		case email == "" && (looksLikeEmailKey(key) || looksLikeEmailValue(value)):
			email = value

		case username == "" && looksLikeUsernameKey(key):
			username = value
		}
	}
	// fallback: if we did not find email by key, try detecting it by value
	if email == "" {
		for _, values := range form {
			if len(values) == 0 {
				continue
			}
			value := strings.TrimSpace(values[0])
			if looksLikeEmailValue(value) {
				email = value
				break
			}
		}
	}
	return email, username, password
}

func normalizeFieldKey(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))

	var b strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func looksLikeEmailKey(key string) bool {
	return strings.Contains(key, "email") || strings.Contains(key, "mail")
}

func looksLikeUsernameKey(key string) bool {
	return strings.Contains(key, "user") ||
		strings.Contains(key, "username") ||
		strings.Contains(key, "login") ||
		strings.Contains(key, "account") ||
		strings.Contains(key, "name")
}

func looksLikePasswordKey(key string) bool {
	return strings.Contains(key, "password") ||
		strings.Contains(key, "passwd") ||
		strings.Contains(key, "pass") ||
		strings.Contains(key, "pwd")
}

func looksLikeEmailValue(value string) bool {
	value = strings.TrimSpace(value)
	return strings.Contains(value, "@") && strings.Contains(value, ".")
}