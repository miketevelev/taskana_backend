package checklists_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

func (s *ChecklistsService) CreateChecklist(
	ctx context.Context,
	userID uuid.UUID,
	checklist domain.Checklist,
) (domain.Checklist, error) {
	if checklist.ID == uuid.Nil {
		checklist.ID = uuid.New()
	}

	if checklist.Version == -1 {
		checklist.Version = 1
	}

	if err := checklist.Validate(); err != nil {
		return domain.Checklist{}, fmt.Errorf(
			"checklist validation failed: %w", err,
		)
	}

	createdChecklist, err := s.checklistRepository.CreateChecklist(
		ctx, userID, checklist,
	)
	if err != nil {
		return domain.Checklist{}, fmt.Errorf(
			"failed to create checklist: %w",
			err,
		)
	}

	return createdChecklist, nil
}
