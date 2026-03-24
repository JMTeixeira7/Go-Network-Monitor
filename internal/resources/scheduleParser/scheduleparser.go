package scheduleParser

import (
	"fmt"
	"strings"
	"time"

	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/model"
)


func ParseScheduleLines(lines []string) ([]*model.Schedule, error) {
	schedules := make([]*model.Schedule, 0, len(lines))

	for i, line := range lines {
		schedule, err := ParseScheduleLine(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", i+1, err)
		}
		schedules = append(schedules, schedule)
	}

	return schedules, nil
}

func ParseScheduleLine(line string) (*model.Schedule, error) {
	fields := strings.Fields(line)
	if len(fields) != 3 {
		return nil, fmt.Errorf("schedule must have exactly 3 fields: <start_time> <end_time> <weekday>")
	}

	startTime, err := ParseClockTime(fields[0])
	if err != nil {
		return nil, fmt.Errorf("parse start time: %w", err)
	}

	endTime, err := ParseClockTime(fields[1])
	if err != nil {
		return nil, fmt.Errorf("parse end time: %w", err)
	}

	weekday, err := ParseWeekday(fields[2])
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
