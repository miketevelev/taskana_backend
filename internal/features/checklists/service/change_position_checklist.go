package checklists_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

func (s *ChecklistsService) ChangePosition(
	ctx context.Context,
	userID uuid.UUID,
	checklistID uuid.UUID,
	newPosition int,
) (domain.Checklist, error) {
	checklist, err := s.GetChecklist(ctx, userID, checklistID)
	if err != nil {
		return domain.Checklist{}, fmt.Errorf(
			"error getting checklist: %w", err,
		)
	}

	oldPosition := checklist.Position

	if oldPosition == newPosition {
		return checklist, nil
	}

	checklist.Position = newPosition

	updatedChecklist, err := s.checklistRepository.ChangePosition(
		ctx, userID, checklist, oldPosition,
	)
	if err != nil {
		return domain.Checklist{}, fmt.Errorf(
			"error changing position in repository: %w", err,
		)
	}

	return updatedChecklist, nil
}
