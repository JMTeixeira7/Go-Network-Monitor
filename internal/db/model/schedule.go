package model

import (
	"time"
)
type Schedule struct {
	Start_time *time.Time
	End_time *time.Time
	Weekday *time.Weekday
}