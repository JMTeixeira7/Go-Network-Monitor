package parser

import (
	"fmt"
	"strings"
	"time"

	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/controller/dto"
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/model"
)

func ParseScheduleLinesToResponses(lines []string) []dto.ScheduleResponse {
	out := make([]dto.ScheduleResponse, 0, len(lines))

	for _, line := range lines {
		req := parseScheduleLine(line)
		out = append(out, dto.ScheduleResponse{
			ID:        "",
			StartTime: req.StartTime,
			EndTime:   req.EndTime,
			Weekday:   req.Weekday,
		})
	}

	return out
}

func ScheduleRequestToScheduleResponses(in []dto.ScheduleRequest) []dto.ScheduleResponse {
	out := make([]dto.ScheduleResponse, 0, len(in))

	for _, s := range in {
		out = append(out, dto.ScheduleResponse{
			ID:        s.ID,
			StartTime: s.StartTime,
			EndTime:   s.EndTime,
			Weekday:   s.Weekday,
		})
	}

	return out
}

func ToDomainSchedules(req dto.BlockedDomainRequest) ([]*model.Schedule, error) {
	if req.SchedulesCount != 0 && req.SchedulesCount != len(req.Schedules) {
		return nil, fmt.Errorf(
			"schedulesCount (%d) does not match schedules length (%d)",
			req.SchedulesCount,
			len(req.Schedules),
		)
	}

	schedules := make([]*model.Schedule, 0, len(req.Schedules))

	for i, s := range req.Schedules {
		domainSchedule, err := ToDomainSchedule(s)
		if err != nil {
			return nil, fmt.Errorf("schedule %d: %w", i+1, err)
		}
		schedules = append(schedules, domainSchedule)
	}

	return schedules, nil
}

func ToDomainSchedule(req dto.ScheduleRequest) (*model.Schedule, error) {
	startTime, err := parseClockString(req.StartTime)
	if err != nil {
		return nil, fmt.Errorf("parse start time: %w", err)
	}

	endTime, err := parseClockString(req.EndTime)
	if err != nil {
		return nil, fmt.Errorf("parse end time: %w", err)
	}

	weekday, err := parseWeekdayString(req.Weekday)
	if err != nil {
		return nil, fmt.Errorf("parse weekday: %w", err)
	}

	schedule, err := model.CreateSchedule(startTime, endTime, weekday)
	if err != nil {
		return nil, fmt.Errorf("create schedule: %w", err)
	}

	return schedule, nil
}

func ToScheduleResponse(s *model.Schedule) dto.ScheduleResponse {
	if s == nil {
		return dto.ScheduleResponse{
			ID:        "",
			StartTime: "",
			EndTime:   "",
			Weekday:   "",
		}
	}

	return dto.ScheduleResponse{
		ID:        "",
		StartTime: formatClock(s.StartTime()),
		EndTime:   formatClock(s.EndTime()),
		Weekday:   formatWeekday(s.Weekday()),
	}
}

func ToScheduleResponses(schedules []*model.Schedule) []dto.ScheduleResponse {
	out := make([]dto.ScheduleResponse, 0, len(schedules))

	for _, s := range schedules {
		out = append(out, ToScheduleResponse(s))
	}

	return out
}

func ToBlockedDomainResponse(domain string, createdAt time.Time, schedules []*model.Schedule) dto.BlockedDomainResponse {
	return dto.BlockedDomainResponse{
		Domain:         domain,
		SchedulesCount: len(schedules),
		CreatedAt:      createdAt.Format(time.RFC3339),
		Schedules:      ToScheduleResponses(schedules),
	}
}

func parseClockString(value string) (*time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" || value == "-" {
		return nil, nil
	}

	layouts := []string{
		"15:04:05",
		"15:04",
	}

	for _, layout := range layouts {
		t, err := time.Parse(layout, value)
		if err == nil {
			return &t, nil
		}
	}

	return nil, fmt.Errorf("invalid clock format %q, expected HH:MM or HH:MM:SS", value)
}

func parseWeekdayString(value string) (*time.Weekday, error) {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" || value == "-" {
		return nil, nil
	}

	var w time.Weekday

	switch value {
	case "sunday", "sun":
		w = time.Sunday
	case "monday", "mon":
		w = time.Monday
	case "tuesday", "tue", "tues":
		w = time.Tuesday
	case "wednesday", "wed":
		w = time.Wednesday
	case "thursday", "thu", "thurs":
		w = time.Thursday
	case "friday", "fri":
		w = time.Friday
	case "saturday", "sat":
		w = time.Saturday
	default:
		return nil, fmt.Errorf("invalid weekday %q", value)
	}

	return &w, nil
}

func formatClock(c *model.Clock) string {
	if c == nil {
		return ""
	}
	return c.String()
}

func formatWeekday(w *time.Weekday) string {
	if w == nil {
		return ""
	}
	return strings.ToLower(w.String())
}