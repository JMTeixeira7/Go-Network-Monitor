package dbmodel

import (
	"fmt"
	"time"
)

type Schedule struct {
	Start_time *Clock
	End_time   *Clock
	Weekday    *time.Weekday
	Timezone   *int
}

type Clock struct {
	hour   int
	min    int
	second int
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

func (c *Clock) String() string {
	if c == nil {
		return "-"
	}
	return fmt.Sprintf("%02d:%02d:%02d", c.hour, c.min, c.second)
}

func ParseClockString(v string) (*Clock, error) {
	t, err := time.Parse("15:04:05", v)
	if err != nil {
		return nil, err
	}

	return &Clock{
		hour:   t.Hour(),
		min:    t.Minute(),
		second: t.Second(),
	}, nil
}