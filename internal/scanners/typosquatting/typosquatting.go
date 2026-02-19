package typosquatting

import (
	"net/http"
)

type Typosquatting struct {
	// add service -_> change design (double dispatch)
}

func New() *Typosquatting {
	return &Typosquatting{}
}

func (t *Typosquatting) Scan(r *http.Request) (res bool, reasons []string) {
	panic("not implemented")
}