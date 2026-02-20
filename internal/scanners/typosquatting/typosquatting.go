package typosquatting

import (
	"context"
	"fmt"
	"net/http"
)

type TyposquattingDBService interface {
	GetVisitedDomains(ctx context.Context) ([]string, error)
	PushDomain(ctx context.Context, domain string) error
}

type Check interface {
	check(reqDomain string, visitedDomain string) bool
	forgeScanMessage(reqDomain string, visitedDomain string) string
}

type additionalCharacter struct{}
type fewerCharacter struct{}
type switchCharacter struct{}

func NewAdditionalCharacterCheck() Check { return additionalCharacter{} }
func NewFewerCharacterCheck() Check      { return fewerCharacter{} }
func NewSwitchCharacterCheck() Check     { return switchCharacter{} }

type Typosquatting struct {
	db_service TyposquattingDBService
	Checks     []Check
}

func New(db_service TyposquattingDBService) *Typosquatting {
	checks := make([]Check, 0, 3)
	checks = append(checks, NewAdditionalCharacterCheck())
	checks = append(checks, NewFewerCharacterCheck())
	checks = append(checks, NewSwitchCharacterCheck()) 
	return &Typosquatting{
		db_service: db_service,
		Checks:     checks}
}

func (t *Typosquatting) Scan(req *http.Request) (res bool, reasons []string) {
	ctx := req.Context()
	visitedDomains, err := t.db_service.GetVisitedDomains(ctx)
	if err != nil {
		fmt.Printf("Error while pulling visited domains from database: %s\n", err)
		return true, nil
	}

	err = t.db_service.PushDomain(ctx, req.URL.Hostname())
	if err != nil {
		fmt.Printf("Error while pushing current request to database: %s\n", err)
		return true, nil
	}

	for i := 0; i < len(visitedDomains); i++ {
		for j := 0; j < len(t.Checks); j++ {
			res = t.Checks[j].check(req.URL.Hostname(), visitedDomains[i])
			if res {
				reasons = append(reasons, t.Checks[j].forgeScanMessage(req.URL.Hostname(), visitedDomains[i]))
				return true, reasons
			}
		}
	}

	return false, reasons
}

func (additionalCharacter) check(reqDomain, visitedDomain string) bool {
	// reqDomain has one extra character compared to visitedDomain
	return oneInsertionAway(visitedDomain, reqDomain)
}

func (fewerCharacter) check(reqDomain, visitedDomain string) bool {
	// reqDomain has one missing character compared to visitedDomain
	return oneInsertionAway(reqDomain, visitedDomain)
}

func (switchCharacter) check(reqDomain, visitedDomain string) bool {
	// swap of two adjacent characters
	return oneAdjacentSwapAway(reqDomain, visitedDomain)
}

func (additionalCharacter) forgeScanMessage(reqDomain, visitedDomain string) string {
	return fmt.Sprintf(
		"Request's domain, %s, has only one additional letter compared with the visited domain, %s.\n", reqDomain, visitedDomain)
}

func (fewerCharacter) forgeScanMessage(reqDomain, visitedDomain string) string {
	return fmt.Sprintf(
		"Request's domain, %s, has only one fewer letter compared with the visited domain, %s.\n", reqDomain, visitedDomain)
}

func (switchCharacter) forgeScanMessage(reqDomain, visitedDomain string) string {
	return fmt.Sprintf(
		"Request's domain, %s, differs only in one letter compared with the visited domain, %s.\n", reqDomain, visitedDomain)
}

// oneInsertionAway returns true if longer can be made equal to shorter
// by removing exactly one character from longer.
// Examples: "google" vs "gooogle" => true
func oneInsertionAway(shorter, longer string) bool {
	if len(longer) != len(shorter)+1 {
		return false
	}

	i, j := 0, 0
	skipped := false

	for i < len(shorter) && j < len(longer) {
		if shorter[i] == longer[j] {
			i++
			j++
			continue
		}
		if skipped {
			return false
		}
		// skip one char in longer
		skipped = true
		j++
	}

	// If we never skipped inside the loop, the "extra" char could be at the end,
	// which still counts as one insertion.
	return true
}

// oneAdjacentSwapAway returns true if a can become b by swapping exactly one
// pair of adjacent characters.
// Example: "googel" vs "google" (swap 'e' and 'l') => true
func oneAdjacentSwapAway(a, b string) bool {
	if len(a) != len(b) || len(a) < 2 {
		return false
	}

	// Find first mismatch
	i := 0
	for i < len(a) && a[i] == b[i] {
		i++
	}
	if i == len(a) {
		return false // identical; not a swap typo
	}

	// i is first mismatch, so we need i+1 to exist for an adjacent swap
	if i+1 >= len(a) {
		return false
	}

	// Check swap at i and i+1
	if a[i] != b[i+1] || a[i+1] != b[i] {
		return false
	}

	// Remaining must match after i+1
	for k := i + 2; k < len(a); k++ {
		if a[k] != b[k] {
			return false
		}
	}

	return true
}
