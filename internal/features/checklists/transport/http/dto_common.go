package checklists_transport_http

import (
	"time"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

type ChecklistDTOResponse struct {
	ID          uuid.UUID `json:"id"`
	Version     int       `json:"version"`
	UserID      uuid.UUID `json:"user_id"`
	TaskID      uuid.UUID `json:"task_id"`
	Title       string    `json:"title"`
	IsCompleted bool      `json:"is_completed"`
	Position    int       `json:"position"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func checklistDTOFromDomain(checklist domain.Checklist) ChecklistDTOResponse {
	return ChecklistDTOResponse{
		ID:          checklist.ID,
		Version:     checklist.Version,
		UserID:      checklist.UserID,
		TaskID:      checklist.TaskID,
		Title:       checklist.Title,
		IsCompleted: checklist.IsCompleted,
		Position:    checklist.Position,
		CreatedAt:   checklist.CreatedAt,
		UpdatedAt:   checklist.UpdatedAt,
	}
}

func checklistDTOsFromDomains(checklists []domain.Checklist) []ChecklistDTOResponse {
	result := make([]ChecklistDTOResponse, len(checklists))
	for i, checklist := range checklists {
		result[i] = checklistDTOFromDomain(checklist)
	}
	return result
}
