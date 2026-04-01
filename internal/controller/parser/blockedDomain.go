package parser

import (
	"fmt"
	"strings"
	"time"

	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/controller/dto"
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/model"
)

// --------------------
// Request DTO -> Domain
// --------------------

func ToDomainSchedules(req dto.BlockedDomainRequest) ([]*model.Schedule, error) {
	if req.SchedulesCount != 0 && req.SchedulesCount != len(req.Schedules) {
		return nil, fmt.Errorf(
			"schedulesCount (%d) does not match schedules length (%d)",
			req.SchedulesCount,
			len(req.Schedules),
		)
	}

	schedules := make([]*model.Schedule, 0, len(req.Schedules))

	for i, scheduleReq := range req.Schedules {
		schedule, err := ToDomainSchedule(scheduleReq)
		if err != nil {
			return nil, fmt.Errorf("schedule %d: %w", i+1, err)
		}
		schedules = append(schedules, schedule)
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

// --------------------
// Domain -> Response DTO
// --------------------

func ToScheduleResponse(schedule *model.Schedule) dto.ScheduleResponse {
	if schedule == nil {
		return dto.ScheduleResponse{
			ID:        "",
			StartTime: "",
			EndTime:   "",
			Weekday:   "",
		}
	}

	return dto.ScheduleResponse{
		ID:        "",
		StartTime: formatClock(schedule.StartTime()),
		EndTime:   formatClock(schedule.EndTime()),
		Weekday:   formatWeekday(schedule.Weekday()),
	}
}

func ToScheduleResponses(schedules []*model.Schedule) []dto.ScheduleResponse {
	out := make([]dto.ScheduleResponse, 0, len(schedules))

	for _, schedule := range schedules {
		out = append(out, ToScheduleResponse(schedule))
	}

	return out
}

func ToBlockedDomainResponse(domain string, createdAt time.Time, schedules []*model.Schedule) dto.BlockedDomainResponse {
	createdAtStr := ""
	if !createdAt.IsZero() {
		createdAtStr = createdAt.Format(time.RFC3339)
	}

	return dto.BlockedDomainResponse{
		Domain:         domain,
		SchedulesCount: len(schedules),
		CreatedAt:      createdAtStr,
		Schedules:      ToScheduleResponses(schedules),
	}
}

// --------------------
// Request DTO -> Response DTO
// Useful when controller wants to echo request data back
// --------------------

func ScheduleRequestToResponse(req dto.ScheduleRequest) dto.ScheduleResponse {
	return dto.ScheduleResponse{
		ID:        req.ID,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Weekday:   req.Weekday,
	}
}

func ScheduleRequestsToResponses(in []dto.ScheduleRequest) []dto.ScheduleResponse {
	out := make([]dto.ScheduleResponse, 0, len(in))

	for _, req := range in {
		out = append(out, ScheduleRequestToResponse(req))
	}

	return out
}

func RequestToBlockedDomainResponse(req dto.BlockedDomainRequest) dto.BlockedDomainResponse {
	return dto.BlockedDomainResponse{
		Domain:         req.Domain,
		SchedulesCount: len(req.Schedules),
		CreatedAt:      req.CreatedAt,
		Schedules:      ScheduleRequestsToResponses(req.Schedules),
	}
}

// --------------------
// Raw schedule lines -> Response DTO
// Useful when domain/service still returns []string
// --------------------

func ScheduleLinesToResponses(lines []string) []dto.ScheduleResponse {
	out := make([]dto.ScheduleResponse, 0, len(lines))

	for _, line := range lines {
		out = append(out, ScheduleLineToResponse(line))
	}

	return out
}

func ScheduleLineToResponse(line string) dto.ScheduleResponse {
	req := parseScheduleLine(line)
	return dto.ScheduleResponse{
		ID:        "",
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Weekday:   req.Weekday,
	}
}

func parseScheduleLine(line string) dto.ScheduleRequest {
	fields := strings.Fields(strings.TrimSpace(line))
	if len(fields) != 3 {
		return dto.ScheduleRequest{
			ID:        "",
			StartTime: "",
			EndTime:   "",
			Weekday:   "",
		}
	}

	startTime := normalizeDashField(fields[0])
	endTime := normalizeDashField(fields[1])
	weekday := normalizeDashField(fields[2])

	return dto.ScheduleRequest{
		ID:        "",
		StartTime: startTime,
		EndTime:   endTime,
		Weekday:   weekday,
	}
}

// --------------------
// Helpers
// --------------------

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

	var weekday time.Weekday

	switch value {
	case "sunday", "sun":
		weekday = time.Sunday
	case "monday", "mon":
		weekday = time.Monday
	case "tuesday", "tue", "tues":
		weekday = time.Tuesday
	case "wednesday", "wed":
		weekday = time.Wednesday
	case "thursday", "thu", "thurs":
		weekday = time.Thursday
	case "friday", "fri":
		weekday = time.Friday
	case "saturday", "sat":
		weekday = time.Saturday
	default:
		return nil, fmt.Errorf("invalid weekday %q", value)
	}

	return &weekday, nil
}

func formatClock(clock *model.Clock) string {
	if clock == nil {
		return ""
	}
	return clock.String()
}

func formatWeekday(weekday *time.Weekday) string {
	if weekday == nil {
		return ""
	}
	return strings.ToLower(weekday.String())
}

func normalizeDashField(value string) string {
	value = strings.TrimSpace(value)
	if value == "-" {
		return ""
	}
	return value
}