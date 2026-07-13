package domain_area_test

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/miketevelev/taskana_backend/internal/core/domain"
	domain_area "github.com/miketevelev/taskana_backend/internal/core/domain/area"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
)

func ptr[T any](v T) *T {
	return &v
}

func TestNewArea(t *testing.T) {
	id := uuid.New()
	userID := uuid.New()
	now := time.Now()

	area := domain_area.NewArea(id, 1, userID, "Work", 2, now, now)

	assert.Equal(t, id, area.ID)
	assert.Equal(t, 1, area.Version)
	assert.Equal(t, userID, area.UserID)
	assert.Equal(t, "Work", area.Title)
	assert.Equal(t, 2, area.Position)
	assert.Equal(t, now, area.CreatedAt)
	assert.Equal(t, now, area.UpdatedAt)
}

func TestNewAreaUninitialized(t *testing.T) {
	userID := uuid.New()
	title := "New Uninitialized Area"
	area := domain_area.NewAreaUninitialized(userID, title)

	assert.Equal(t, domain.UninitializedID, area.ID)
	assert.Equal(t, domain.UninitializedVersion, area.Version)
	assert.Equal(t, userID, area.UserID)
	assert.Equal(t, title, area.Title)
	assert.Equal(t, 1, area.Position)
	assert.NotZero(t, area.CreatedAt)
	assert.NotZero(t, area.UpdatedAt)
	assert.Equal(t, area.CreatedAt, area.UpdatedAt)
}

func TestArea_Validate(t *testing.T) {
	now := time.Now()
	validArea := domain_area.NewArea(
		uuid.New(), 1, uuid.New(), "Work", 1, now, now,
	)

	tests := []struct {
		name        string
		mutateArea  func(a *domain_area.Area)
		expectedErr error
		checkResult func(t *testing.T, a *domain_area.Area)
	}{
		{
			name:        "valid area",
			mutateArea:  func(a *domain_area.Area) {},
			expectedErr: nil,
		},
		{
			name: "title too short",
			mutateArea: func(a *domain_area.Area) {
				a.Title = "Hi"
			},
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "title only spaces (fails trim check)",
			mutateArea: func(a *domain_area.Area) {
				a.Title = "     "
			},
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "title too long",
			mutateArea: func(a *domain_area.Area) {
				a.Title = strings.Repeat("a", 101)
			},
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "position less than 1",
			mutateArea: func(a *domain_area.Area) {
				a.Position = 0
			},
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "created_at is zero",
			mutateArea: func(a *domain_area.Area) {
				a.CreatedAt = time.Time{}
			},
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "updated_at is zero",
			mutateArea: func(a *domain_area.Area) {
				a.UpdatedAt = time.Time{}
			},
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "updated before created",
			mutateArea: func(a *domain_area.Area) {
				a.UpdatedAt = a.CreatedAt.Add(-1 * time.Hour)
			},
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "negative version gets mutated to 0",
			mutateArea: func(a *domain_area.Area) {
				a.Version = -5
			},
			expectedErr: nil,
			checkResult: func(t *testing.T, a *domain_area.Area) {
				assert.Equal(t, 0, a.Version)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				area := validArea
				tt.mutateArea(&area)

				err := area.Validate()

				if tt.expectedErr != nil {
					require.Error(t, err)
					assert.ErrorIs(t, err, tt.expectedErr)
				} else {
					assert.NoError(t, err)
				}

				if tt.checkResult != nil {
					tt.checkResult(t, &area)
				}
			},
		)
	}
}

func TestNewAreaPatch(t *testing.T) {
	titleNullable := domain.Nullable[string]{
		Set: true, Value: ptr("Patched Title"),
	}
	patch := domain_area.NewAreaPatch(titleNullable)

	assert.Equal(t, titleNullable, patch.Title)
}

func TestAreaPatch_Validate(t *testing.T) {
	tests := []struct {
		name        string
		patch       domain_area.AreaPatch
		expectedErr error
	}{
		{
			name: "valid empty patch",
			patch: domain_area.AreaPatch{
				Title: domain.Nullable[string]{Set: false, Value: nil},
			},
			expectedErr: nil,
		},
		{
			name: "valid title patch",
			patch: domain_area.AreaPatch{
				Title: domain.Nullable[string]{
					Set: true, Value: ptr("New Title"),
				},
			},
			expectedErr: nil,
		},
		{
			name: "invalid patch: Set is true but Value is nil",
			patch: domain_area.AreaPatch{
				Title: domain.Nullable[string]{Set: true, Value: nil},
			},
			expectedErr: core_errors.ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				err := tt.patch.Validate()

				if tt.expectedErr != nil {
					require.Error(t, err)
					assert.ErrorIs(t, err, tt.expectedErr)
				} else {
					assert.NoError(t, err)
				}
			},
		)
	}
}

func TestArea_ApplyPatch(t *testing.T) {
	now := time.Now()
	validArea := domain_area.NewArea(
		uuid.New(), 1, uuid.New(), "Initial Title", 1, now, now,
	)

	tests := []struct {
		name          string
		patch         domain_area.AreaPatch
		expectedErr   error
		expectedTitle string
	}{
		{
			name: "successfully apply title patch",
			patch: domain_area.AreaPatch{
				Title: domain.Nullable[string]{
					Set: true, Value: ptr("Updated Title"),
				},
			},
			expectedErr:   nil,
			expectedTitle: "Updated Title",
		},
		{
			name: "successfully apply empty patch (no changes)",
			patch: domain_area.AreaPatch{
				Title: domain.Nullable[string]{Set: false, Value: nil},
			},
			expectedErr:   nil,
			expectedTitle: "Initial Title",
		},
		{
			name: "fail: patch is invalid",
			patch: domain_area.AreaPatch{
				Title: domain.Nullable[string]{Set: true, Value: nil},
			},
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "fail: patched area becomes invalid (title too short)",
			patch: domain_area.AreaPatch{
				Title: domain.Nullable[string]{Set: true, Value: ptr("No")},
			},
			expectedErr: core_errors.ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				area := validArea
				err := area.ApplyPatch(tt.patch)

				if tt.expectedErr != nil {
					require.Error(t, err)
					assert.ErrorIs(t, err, tt.expectedErr)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tt.expectedTitle, area.Title)
				}
			},
		)
	}
}
