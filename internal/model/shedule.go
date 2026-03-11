package model

import (
	"fmt"
	"time"
)

type Schedule struct {

	start_time *time.Time
	end_time *time.Time
	weekday *time.Weekday
}

func CreateSchedule(start_time *time.Time, end_time *time.Time, weekday *time.Weekday) (*Schedule, error) {
	if start_time == nil && end_time == nil && weekday == nil {
		return nil, nil
	}
	if (start_time != nil && end_time == nil) || (start_time == nil && end_time != nil) {
		return nil, fmt.Errorf("schedule format is wrong: start_time and end_time must both be set or both be nil")
	}
	if start_time != nil && end_time != nil && end_time.Before(*start_time) {
		return nil, fmt.Errorf("schedule format is wrong: end_time is before start_time")
	}
	return &Schedule{
		start_time: start_time,
		end_time: end_time,
		weekday: weekday,
	}, nil
}

func (s *Schedule) StartTime() *time.Time {
	if s == nil {
		return nil
	}
	return s.start_time
}

func (s *Schedule) EndTime() *time.Time {
	if s == nil {
		return nil
	}
	return s.end_time
}

func (s *Schedule) Weekday() *time.Weekday {
	if s == nil {
		return nil
	}
	return s.weekday
}