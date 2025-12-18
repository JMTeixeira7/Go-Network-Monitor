package url_error

import (
	"fmt"
)
type TargetUrlError struct {
	url string
	domain string
}

func NewTargetUrlError(url, domain string) error {
	return &TargetUrlError{url: url, domain: domain}
}

func (e *TargetUrlError) Error() string {
	return fmt.Sprintf("Access to this domain, %s, is restricted. Cannot access the requested URL, %s", e.domain, e.url)
}