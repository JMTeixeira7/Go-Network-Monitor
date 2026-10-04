package model

import (
	"time"

	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/db/dbmodel"
)

type Domain struct {
	Domain string
	Time time.Time
}

func CreateDomain(domain string, schedule time.Time) Domain {
	return Domain{
		Domain: domain,
		Time: schedule,
	}
}

func CreateDomainsFromDBDomains(domains_db []dbmodel.Domain) []Domain {
	if domains_db == nil {
		return nil
	}

	var domains_model []Domain
	for _, domain_db := range domains_db{
		domain_model := CreateDomain(domain_db.Domain, domain_db.Time)
		domains_model = append(domains_model, domain_model)
	}

	return domains_model 
}
