package model

import (
	"fmt"
	"time"
)

type Schedule struct {
	start_time *Clock
	end_time   *Clock
	weekday    *time.Weekday
	timezone   *int
}

type Clock struct {
	hour   int
	min    int
	second int
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

	_, offset := time.Now().Zone()

	s := &Schedule{
		weekday:  weekday,
		timezone: &offset,
	}

	if start_time != nil {
		s.start_time = &Clock{
			hour:   start_time.Hour(),
			min:    start_time.Minute(),
			second: start_time.Second(),
		}
	}

	if end_time != nil {
		s.end_time = &Clock{
			hour:   end_time.Hour(),
			min:    end_time.Minute(),
			second: end_time.Second(),
		}
	}

	return s, nil
}

func CreateScheduleFromDB(start_time *Clock, end_time *Clock, weekday *time.Weekday, timezone *int) (*Schedule, error) {
	if start_time == nil && end_time == nil && weekday == nil && timezone == nil {
		return nil, nil
	}
	return &Schedule{
		start_time: start_time,
		end_time: end_time,
		weekday: weekday,
		timezone: timezone,
	}, nil
}

func (s *Schedule) StartTime() *Clock {
	if s == nil || s.start_time == nil {
		return nil
	}
	return s.start_time
}

func (s *Schedule) EndTime() *Clock {
	if s == nil || s.end_time == nil {
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

func (s *Schedule) Timezone() *int {
	if s == nil {
		return nil
	}
	return s.timezone
}

func (c *Clock) String() string {
	if c == nil {
		return "-"
	}
	return fmt.Sprintf("%02d:%02d:%02d", c.hour, c.min, c.second)
}

func CreateClock(hour int, min int, seconds int) *Clock {
	return &Clock{
		hour: hour,
		min: min,
		second: seconds,
	}
}

func (c *Clock) GetHour() int {
	return c.hour
}

func (c *Clock) GetMin() int {
	return c.min
}

func (c *Clock) GetSeconds() int {
	return c.second
}
