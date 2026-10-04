package dbmodel

import "time"

type Domain struct {
	ID	int64
	Domain string
	Time time.Time
}