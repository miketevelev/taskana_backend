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
	GetChecklist(
		ctx context.Context,
		userID uuid.UUID,
		checklistID uuid.UUID,
	) (domain.Checklist, error)

	GetChecklists(
		ctx context.Context,
		userID uuid.UUID,
		limit *int,
		offset *int,
	) ([]domain.Checklist, error)

	CreateChecklist(
		ctx context.Context,
		userID uuid.UUID,
		checklist domain.Checklist,
	) (domain.Checklist, error)

	ChangePosition(
		ctx context.Context,
		userID uuid.UUID,
		checklist domain.Checklist,
		oldPosition int,
	) (domain.Checklist, error)

	PatchChecklist(
		ctx context.Context,
		userID uuid.UUID,
		checklist domain.Checklist,
	) (domain.Checklist, error)

	DeleteChecklist(
		ctx context.Context,
		userID uuid.UUID,
		checklistID uuid.UUID,
	) error
}

func NewChecklistsService(
	checklistRepository ChecklistsRepository,
) *ChecklistsService {
	return &ChecklistsService{
		checklistRepository: checklistRepository,
	}
}
