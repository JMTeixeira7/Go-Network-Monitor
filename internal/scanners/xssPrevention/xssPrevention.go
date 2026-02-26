package xssprevention

import (
	"fmt"
	"net/http"

	"github.com/microcosm-cc/bluemonday"
)

type XSSPrevent struct {
	policy *bluemonday.Policy
	maxPreview int64
}

func New() *XSSPrevent {
	return &XSSPrevent{
		policy: bluemonday.StrictPolicy(),
	}
}

func (x *XSSPrevent) Scan(r *http.Request) (bool, []string) {
	q := r.URL.Query()
	for key, vals := range q {
		for _, v := range vals {
			if x.policy.Sanitize(v) != v {
				return true, []string{fmt.Sprintf("blocked: query parameter: \"%s\", modified by sanitizer.", key)}
			}
		}
	}
	return false, nil
}