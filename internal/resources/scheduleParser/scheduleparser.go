package scheduleParser

import (
	"fmt"
	"strings"
	"time"

	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/controller"
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/model"
)


func ParseBlockedDomainSchedules(input controller.BlockedDomain) ([]*model.Schedule, error) {
	if input.SchedulesCount != 0 && input.SchedulesCount != len(input.Schedules) {
		return nil, fmt.Errorf(
			"schedulesCount (%d) does not match schedules length (%d)",
			input.SchedulesCount,
			len(input.Schedules),
		)
	}

	return ParseSchedules(input.Schedules)
}

func ParseSchedules(items []controller.Schedule) ([]*model.Schedule, error) {
	schedules := make([]*model.Schedule, 0, len(items))

	for i, item := range items {
		schedule, err := ParseSchedule(item)
		if err != nil {
			return nil, fmt.Errorf("schedule %d: %w", i+1, err)
		}
		schedules = append(schedules, schedule)
	}

	return schedules, nil
}

func ParseSchedule(item controller.Schedule) (*model.Schedule, error) {
	// Optional: treat an empty schedule as nil, similar to "- - -"
	if strings.TrimSpace(item.StartTime) == "" &&
		strings.TrimSpace(item.EndTime) == "" &&
		strings.TrimSpace(item.Weekday) == "" {
		return nil, nil
	}

	startTime, err := ParseClockTime(item.StartTime)
	if err != nil {
		return nil, fmt.Errorf("parse start time: %w", err)
	}

	endTime, err := ParseClockTime(item.EndTime)
	if err != nil {
		return nil, fmt.Errorf("parse end time: %w", err)
	}

	weekday, err := ParseWeekday(item.Weekday)
	if err != nil {
		return nil, fmt.Errorf("parse weekday: %w", err)
	}

	schedule, err := model.CreateSchedule(startTime, endTime, weekday)
	if err != nil {
		return nil, fmt.Errorf("create schedule: %w", err)
	}

	return schedule, nil
}

func FormatSchedules(schedules []*model.Schedule) ([]string, error) {
	lines := make([]string, 0, len(schedules))

	for _, schedule := range schedules {
		if schedule == nil {
			lines = append(lines, "- - -")
			continue
		}

		start := schedule.StartTime()
		end := schedule.EndTime()

		weekday := "-"
		if w := schedule.Weekday(); w != nil {
			weekday = WeekdayToString(*w)
		}

		lines = append(lines, fmt.Sprintf("%s %s %s", start, end, weekday))
	}

	return lines, nil
}

func ParseClockTime(s string) (*time.Time, error) {
	if s == "-" || s == "" {
		return nil, nil
	}

	for _, layout := range []string{"15:04:05", "15:04"} {
		t, err := time.Parse(layout, s)
		if err == nil {
			return &t, nil
		}
	}

	return nil, fmt.Errorf("invalid clock time %q: expected HH:MM or HH:MM:SS", s)
}

func ParseWeekday(s string) (*time.Weekday, error) {
	if s == "-" || s == "" {
		return nil, nil
	}

	switch strings.ToLower(s) {
	case "sunday":
		w := time.Sunday
		return &w, nil
	case "monday":
		w := time.Monday
		return &w, nil
	case "tuesday":
		w := time.Tuesday
		return &w, nil
	case "wednesday":
		w := time.Wednesday
		return &w, nil
	case "thursday":
		w := time.Thursday
		return &w, nil
	case "friday":
		w := time.Friday
		return &w, nil
	case "saturday":
		w := time.Saturday
		return &w, nil
	default:
		return nil, fmt.Errorf("invalid weekday %q", s)
	}
}

func WeekdayToString(w time.Weekday) string {
	switch w {
	case time.Sunday:
		return "Sunday"
	case time.Monday:
		return "Monday"
	case time.Tuesday:
		return "Tuesday"
	case time.Wednesday:
		return "Wednesday"
	case time.Thursday:
		return "Thursday"
	case time.Friday:
		return "Friday"
	case time.Saturday:
		return "Saturday"
	default:
		return "-"
	}
}
