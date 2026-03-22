package dbmodel

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
	minute int
	second int
}

func NewClock(hour, minute, second int) (*Clock, error) {
	if hour < 0 || hour > 23 {
		return nil, fmt.Errorf("invalid hour %d", hour)
	}
	if minute < 0 || minute > 59 {
		return nil, fmt.Errorf("invalid minute %d", minute)
	}
	if second < 0 || second > 59 {
		return nil, fmt.Errorf("invalid second %d", second)
	}

	return &Clock{
		hour:   hour,
		minute: minute,
		second: second,
	}, nil
}

func ParseClockString(v string) (*Clock, error) {
	t, err := time.Parse("15:04:05", v)
	if err != nil {
		return nil, fmt.Errorf("parse clock %q: %w", v, err)
	}

	return NewClock(t.Hour(), t.Minute(), t.Second())
}

func NewSchedule(startTime, endTime *Clock, weekday *time.Weekday, timezone *int) (*Schedule, error) {
	// Represents "no schedule row semantics".
	if startTime == nil && endTime == nil && weekday == nil && timezone == nil {
		return nil, nil
	}

	// Weekday-only schedule is allowed.
	if startTime == nil && endTime == nil {
		if timezone != nil {
			return nil, fmt.Errorf("timezone requires start and end time")
		}
		return &Schedule{
			weekday: weekday,
		}, nil
	}

	// From here on, time windows must be complete and timezone-aware.
	if startTime == nil || endTime == nil {
		return nil, fmt.Errorf("start and end time must both be set")
	}
	if timezone == nil {
		return nil, fmt.Errorf("timezone must be set when start and end time are set")
	}
	if !startTime.Before(endTime) {
		return nil, fmt.Errorf("end time must be after start time")
	}

	return &Schedule{
		startTime: startTime,
		endTime:   endTime,
		weekday:   weekday,
		timezone:  timezone,
	}, nil
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

func (c *Clock) Hour() int {
	return c.hour
}

func (c *Clock) Minute() int {
	return c.minute
}

func (c *Clock) Second() int {
	return c.second
}

func (c *Clock) Seconds() int {
	if c == nil {
		return 0
	}
	return c.hour*3600 + c.minute*60 + c.second
}

func (c *Clock) Before(other *Clock) bool {
	if c == nil || other == nil {
		return false
	}
	return c.Seconds() < other.Seconds()
}

func (c *Clock) String() string {
	if c == nil {
		return ""
	}
	return fmt.Sprintf("%02d:%02d:%02d", c.hour, c.minute, c.second)
}

func (s *Schedule) Matches(now *time.Time, day *time.Weekday) bool {
	// nil schedule means "always blocked"
	if s == nil {
		return true
	}

	// If schedule is weekday-scoped, weekday must match first.
	if s.weekday != nil {
		if day == nil || *day != *s.weekday {
			return false
		}
	}

	// If there is no time window, weekday-only rule applies.
	if s.startTime == nil || s.endTime == nil || now == nil {
		return true
	}

	nowSec := now.Hour()*3600 + now.Minute()*60 + now.Second()
	return nowSec >= s.startTime.Seconds() && nowSec < s.endTime.Seconds()
}

func AnyScheduleMatches(schedules []*Schedule, now *time.Time, day *time.Weekday) bool {
	for _, s := range schedules {
		if s.Matches(now, day) {
			return true
		}
	}
	return false
}

// SQLValues centralizes how a dbmodel.Schedule is stored in SQL.
func (s *Schedule) SQLValues() (start any, end any, weekday any, timezone any) {
	if s == nil {
		return nil, nil, nil, nil
	}

	if s.startTime != nil {
		start = s.startTime.String()
	}
	if s.endTime != nil {
		end = s.endTime.String()
	}
	if s.weekday != nil {
		weekday = int(*s.weekday)
	}
	if s.timezone != nil {
		timezone = *s.timezone
	}

	return start, end, weekday, timezone
}
