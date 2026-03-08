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
		return &Schedule{
				start_time: nil,
				end_time: nil,
				weekday: nil,
			}, fmt.Errorf("Schedule format is wrong.\n")
	}
	return &Schedule{
		start_time: start_time,
		end_time: end_time,
		weekday: weekday,
	}, nil
}