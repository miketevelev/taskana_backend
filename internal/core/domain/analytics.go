package domain

import (
	"time"

	"github.com/google/uuid"
)

type AnalyticsResponse struct {
	TimeByProject      []ProjectTimeSlice      `json:"time_by_project"`
	FocusTimeline      []DailyFocusSlice       `json:"focus_timeline"`
	EstimationAccuracy []EstimationAccuracyRow `json:"estimation_accuracy"`
}

type ProjectTimeSlice struct {
	ProjectID    uuid.UUID `json:"project_id"`
	ProjectTitle string    `json:"project_title"`
	TotalSeconds int64     `json:"total_seconds"`
}

type DailyFocusSlice struct {
	FocusDay      time.Time `json:"focus_day"`
	DailySeconds  int64     `json:"daily_seconds"`
	SessionsCount int64     `json:"sessions_count"`
}

type EstimationAccuracyRow struct {
	TaskID             uuid.UUID `json:"task_id"`
	TaskTitle          string    `json:"task_title"`
	PlannedPomodoros   int       `json:"planned_pomodoros"`
	ActualPomodoros    int64     `json:"actual_pomodoros"`
	AccuracyPercentage float64   `json:"accuracy_percentage"`
}
