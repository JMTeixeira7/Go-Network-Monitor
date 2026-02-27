package model

import "github.com/JMTeixeira7/Go-Network-Monitor.git/internal/security"


type Credentials struct {
	Username string
	Fingerprint string
}

func CreateCredentials(email, username, password string, fp *security.Fingerprinter) *Credentials {
	u := username
	if u == "" {
		u = email
	}
	if u == "" || password == "" {
		return nil
	}
	return &Credentials{
		Username:    u,
		Fingerprint: fp.FingerprintPasswordHex(password),
	}
}


