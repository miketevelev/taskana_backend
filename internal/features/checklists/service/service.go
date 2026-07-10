package checklists_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

type ChecklistsService struct {
	checklistRepository ChecklistsRepository
}

type ChecklistsRepository interface {
	CreateChecklist(
		ctx context.Context,
		userID uuid.UUID,
		checklist domain.Checklist,
	) (domain.Checklist, error)
}

func NewChecklistsService(
	checklistRepository ChecklistsRepository,
) *ChecklistsService {
	return &ChecklistsService{
		checklistRepository: checklistRepository,
	}
}
