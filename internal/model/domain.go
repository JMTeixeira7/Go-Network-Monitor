package model

import "time"

type Domain struct {
	Domain string
	Time time.Time
}

func CreateDomain(domain *string, schedule *time.Time) *Domain {
	if domain == nil || schedule == nil {
		return nil
	}

	return &Domain{
		Domain: *domain,
		Time: *schedule,
	}
}
