package dto

type VisitedDomainRequest struct {
	Offset string `json:"offset"`
	Limit string `json:"limit"`
}

type VisitedDomainResponse struct {
	VisitedDomains []Domain `json:"VisitedDomains"`
}

type Domain struct{
	Domain string `json:"Domain"`
	Time string `json:"Time"` 
}