package storage

import (
	"time"
)

type Search struct {
	server string
	method string
	domain string
	time time.Time //date too

}