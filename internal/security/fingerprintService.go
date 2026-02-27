package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

type Fingerprinter struct {
	seed []byte
}

func NewFingerprinter(seed []byte) *Fingerprinter {
	return &Fingerprinter{seed: seed}
}

func (f *Fingerprinter) FingerprintPasswordHex(password string) string {
	mac := hmac.New(sha256.New, f.seed)
	mac.Write([]byte(password))
	return hex.EncodeToString(mac.Sum(nil))
}
