package credentialsParser

import (
	"net/url"
	"strings"
	"unicode"
)

func ExtractCredentialFields(form url.Values) (email, username, password string) {
	for rawKey, values := range form {
		if len(values) == 0 {
			continue
		}

		value := strings.TrimSpace(values[0])
		if value == "" {
			continue
		}

		key := NormalizeFieldKey(rawKey)

		switch {
		case password == "" && LooksLikePasswordKey(key):
			password = value

		case email == "" && (LooksLikeEmailKey(key) || LooksLikeEmailValue(value)):
			email = value

		case username == "" && LooksLikeUsernameKey(key):
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
			if LooksLikeEmailValue(value) {
				email = value
				break
			}
		}
	}
	return email, username, password
}

func NormalizeFieldKey(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))

	var b strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func LooksLikeEmailKey(key string) bool {
	return strings.Contains(key, "email") || strings.Contains(key, "mail")
}

func LooksLikeUsernameKey(key string) bool {
	return strings.Contains(key, "user") ||
		strings.Contains(key, "username") ||
		strings.Contains(key, "login") ||
		strings.Contains(key, "account") ||
		strings.Contains(key, "name")
}

func LooksLikePasswordKey(key string) bool {
	return strings.Contains(key, "password") ||
		strings.Contains(key, "passwd") ||
		strings.Contains(key, "pass") ||
		strings.Contains(key, "pwd")
}

func LooksLikeEmailValue(value string) bool {
	value = strings.TrimSpace(value)
	return strings.Contains(value, "@") && strings.Contains(value, ".")
}
