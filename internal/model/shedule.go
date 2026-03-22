package model

import (
	"fmt"
	"time"
)

type Schedule struct {
	startTime *Clock
	endTime   *Clock
	weekday   *time.Weekday
	timezone  *int
}

type Clock struct {
	hour   int
	min    int
	second int
}

func CreateSchedule(startTime *time.Time, endTime *time.Time, weekday *time.Weekday) (*Schedule, error) {
	if startTime == nil && endTime == nil && weekday == nil {
		return nil, nil
	}

	_, offset := time.Now().Zone()

	s := &Schedule{
		weekday:  weekday,
		timezone: &offset,
	}

	if startTime != nil {
		s.startTime = clockFromTime(startTime)
	}

	if endTime != nil {
		s.endTime = clockFromTime(endTime)
	}

	return s, nil
}

func CreateScheduleFromDB(startTime *Clock, endTime *Clock, weekday *time.Weekday, timezone *int) (*Schedule, error) {
	if startTime == nil && endTime == nil && weekday == nil && timezone == nil {
		return nil, nil
	}

	return &Schedule{
		startTime: startTime,
		endTime:   endTime,
		weekday:   weekday,
		timezone:  timezone,
	}, nil
}

func clockFromTime(t *time.Time) *Clock {
	if t == nil {
		return nil
	}

	return &Clock{
		hour:   t.Hour(),
		min:    t.Minute(),
		second: t.Second(),
	}
}

func (s *Schedule) StartTime() *Clock {
	if s == nil {
		return nil
	}
	return s.startTime
}

func (s *Schedule) EndTime() *Clock {
	if s == nil {
		return nil
	}
	return s.endTime
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

func CreateClock(hour int, min int, second int) *Clock {
	return &Clock{
		hour:   hour,
		min:    min,
		second: second,
	}
}

func (c *Clock) Hour() int {
	return c.hour
}

func (c *Clock) Minute() int {
	return c.min
}

func (c *Clock) Second() int {
	return c.second
}
