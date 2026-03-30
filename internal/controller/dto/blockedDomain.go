package dto


type BlockedDomainRequest struct {
	Domain         string     `json:"domain"`
	SchedulesCount int        `json:"schedulesCount"`
	CreatedAt      string     `json:"createdAt"`
	Schedules      []ScheduleRequest `json:"schedules"`
}

type ScheduleRequest struct {
	ID        string `json:"id"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Weekday   string `json:"weekday"`
}

type BlockedDomainResponse struct {
	Domain         string     `json:"domain"`
	SchedulesCount int        `json:"schedulesCount"`
	CreatedAt      string     `json:"createdAt"`
	Schedules      []ScheduleResponse `json:"schedules"`
}

type ScheduleResponse struct {
	ID        string `json:"id"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Weekday   string `json:"weekday"`
}