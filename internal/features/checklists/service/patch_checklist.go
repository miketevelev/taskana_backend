package checklists_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

func (s *ChecklistsService) PatchChecklist(
	ctx context.Context,
	userID uuid.UUID,
	checklistID uuid.UUID,
	patch domain.ChecklistPatch,
) (domain.Checklist, error) {
	checklist, err := s.checklistRepository.GetChecklist(
		ctx, userID, checklistID,
	)
	if err != nil {
		return domain.Checklist{}, fmt.Errorf(
			"error while fetching checklist: %w", err,
		)
	}

	if err := checklist.ApplyPatch(patch); err != nil {
		return domain.Checklist{}, fmt.Errorf(
			"error while applying patch to checklist: %w", err,
		)
	}

	patchedChecklist, err := s.checklistRepository.PatchChecklist(
		ctx, userID, checklist,
	)
	if err != nil {
		return domain.Checklist{}, fmt.Errorf(
			"error while saving patched checklist: %w", err,
		)
	}

	return patchedChecklist, nil
}
