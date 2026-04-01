package dto

type ApiResponse[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type StatusResponse struct {
	ListenerRunning bool   `json:"listenerRunning"`
	CacheStatus     string `json:"cacheStatus"`
	LastUpdated     string `json:"lastUpdated"`
}

type ListenerStateResponse struct {
	ListenerRunning bool `json:"listenerRunning"`
}

type CacheStateResponse struct {
	CacheCleared bool `json:"cacheCleared"`
}

type BlockedDomainRequest struct {
	Domain         string            `json:"domain"`
	SchedulesCount int               `json:"schedulesCount"`
	CreatedAt      string            `json:"createdAt"`
	Schedules      []ScheduleRequest `json:"schedules"`
}

type ScheduleRequest struct {
	ID        string `json:"id"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Weekday   string `json:"weekday"`
}

type BlockedDomainResponse struct {
	Domain         string             `json:"domain"`
	SchedulesCount int                `json:"schedulesCount"`
	CreatedAt      string             `json:"createdAt"`
	Schedules      []ScheduleResponse `json:"schedules"`
}

type ScheduleResponse struct {
	ID        string `json:"id"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Weekday   string `json:"weekday"`
}