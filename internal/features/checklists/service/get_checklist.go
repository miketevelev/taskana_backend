package checklists_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

func (s *ChecklistsService) GetChecklist(
	ctx context.Context,
	userID uuid.UUID,
	checklistID uuid.UUID,
) (domain.Checklist, error) {
	checklist, err := s.checklistRepository.GetChecklist(
		ctx, userID, checklistID,
	)
	if err != nil {
		return domain.Checklist{}, fmt.Errorf(
			"error getting checklist: %w", err,
		)
	}

	return checklist, nil
}
