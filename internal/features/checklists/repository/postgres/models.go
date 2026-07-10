package checklists_postgres_repository

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

type ChecklistModel struct {
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

func checklistDomainFromModel(checklistModel ChecklistModel) domain.Checklist {
	return domain.Checklist{
		ID:          checklistModel.ID,
		Version:     checklistModel.Version,
		UserID:      checklistModel.UserID,
		TaskID:      checklistModel.TaskID,
		Title:       checklistModel.Title,
		IsCompleted: checklistModel.IsCompleted,
		Position:    checklistModel.Position,
		CreatedAt:   checklistModel.CreatedAt,
		UpdatedAt:   checklistModel.UpdatedAt,
	}
}

func checklistDomainsFromModels(checklists []ChecklistModel) []domain.Checklist {
	if len(checklists) == 0 {
		return []domain.Checklist{}
	}
	checklistDomains := make([]domain.Checklist, len(checklists))

	for i, checklist := range checklists {
		checklistDomains[i] = domain.Checklist{
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

	return checklistDomains
}

func scanChecklist(row interface{ Scan(dest ...any) error }) (
	ChecklistModel, error,
) {
	var checklistModel ChecklistModel
	err := row.Scan(
		&checklistModel.ID,
		&checklistModel.Version,
		&checklistModel.UserID,
		&checklistModel.TaskID,
		&checklistModel.Title,
		&checklistModel.IsCompleted,
		&checklistModel.Position,
		&checklistModel.CreatedAt,
		&checklistModel.UpdatedAt,
	)
	if err != nil {
		return ChecklistModel{}, fmt.Errorf("scan checklist: %w", err)
	}
	return checklistModel, nil
}
