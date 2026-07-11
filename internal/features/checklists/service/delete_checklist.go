package checklists_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *ChecklistsService) DeleteChecklist(
	ctx context.Context,
	userID uuid.UUID,
	checklistID uuid.UUID,
) error {
	if err := s.checklistRepository.DeleteChecklist(
		ctx, userID, checklistID,
	); err != nil {
		return fmt.Errorf(
			"delete checklist: %w", err,
		)
	}

	return nil
}
