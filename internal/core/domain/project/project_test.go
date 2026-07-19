package domain_project_test

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	domain_project "github.com/miketevelev/taskana_backend/internal/core/domain/project"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
)

func ptr[T any](v T) *T {
	return &v
}

// nullableSet builds a domain.Nullable[T] that is "set" to a concrete value.
// For pointer-typed T (e.g. *uuid.UUID, *string, *time.Time) this produces a
// Nullable whose Value is a pointer-to-pointer, matching how ProjectPatch
// fields are declared (e.g. domain.Nullable[*uuid.UUID]).
func nullableSet[T any](v T) domain.Nullable[T] {
	return domain.Nullable[T]{Set: true, Value: ptr(v)}
}

// nullableUnset builds a domain.Nullable[T] that was never touched by the patch.
func nullableUnset[T any]() domain.Nullable[T] {
	return domain.Nullable[T]{Set: false, Value: nil}
}

// nullableClear builds a domain.Nullable[T] that is explicitly set to NULL —
// i.e. Set is true but Value itself is nil. For scalar fields like Title/Status
// this is an invalid patch (can't null a required field). For pointer-typed
// fields like AreaID/Notes/Deadline/CompletedAt this is how callers clear the
// field (e.g. detach a project from its area).
func nullableClear[T any]() domain.Nullable[T] {
	return domain.Nullable[T]{Set: true, Value: nil}
}

func TestNewProject(t *testing.T) {
	id := uuid.New()
	userID := uuid.New()
	areaID := uuid.New()
	now := time.Now()
	deadline := now.Add(48 * time.Hour)
	completedAt := now.Add(time.Hour)
	notes := "some notes"

	p := domain_project.NewProject(
		id, 3, userID, &areaID, "Launch", &notes,
		domain_project.ProjectStatusCompleted, 2, &deadline, &completedAt, now,
		now,
	)

	assert.Equal(t, id, p.ID)
	assert.Equal(t, 3, p.Version)
	assert.Equal(t, userID, p.UserID)
	assert.Equal(t, &areaID, p.AreaID)
	assert.Equal(t, "Launch", p.Title)
	assert.Equal(t, &notes, p.Notes)
	assert.Equal(t, domain_project.ProjectStatusCompleted, p.Status)
	assert.Equal(t, 2, p.Position)
	assert.Equal(t, &deadline, p.Deadline)
	assert.Equal(t, &completedAt, p.CompletedAt)
	assert.Equal(t, now, p.CreatedAt)
	assert.Equal(t, now, p.UpdatedAt)
}

func TestNewProjectUninitialized(t *testing.T) {
	userID := uuid.New()
	areaID := uuid.New()
	notes := "notes"
	deadline := time.Now().Add(24 * time.Hour)

	p := domain_project.NewProjectUninitialized(
		userID, &areaID, "New Project", &notes, &deadline,
	)

	assert.Equal(t, domain.UninitializedID, p.ID)
	assert.Equal(t, domain.UninitializedVersion, p.Version)
	assert.Equal(t, userID, p.UserID)
	assert.Equal(t, &areaID, p.AreaID)
	assert.Equal(t, "New Project", p.Title)
	assert.Equal(t, &notes, p.Notes)
	assert.Equal(t, domain_project.ProjectStatusActive, p.Status)
	assert.Equal(t, 1, p.Position)
	assert.Equal(t, &deadline, p.Deadline)
	assert.Nil(t, p.CompletedAt)
	assert.NotZero(t, p.CreatedAt)
	assert.NotZero(t, p.UpdatedAt)
	assert.Equal(t, p.CreatedAt, p.UpdatedAt)
}

func TestNewProjectUninitialized_NilOptionalFields(t *testing.T) {
	userID := uuid.New()

	p := domain_project.NewProjectUninitialized(
		userID, nil, "New Project", nil, nil,
	)

	assert.Nil(t, p.AreaID)
	assert.Nil(t, p.Notes)
	assert.Nil(t, p.Deadline)
	assert.Nil(t, p.CompletedAt)
}

func validProject(now time.Time) domain_project.Project {
	return domain_project.NewProject(
		uuid.New(), 1, uuid.New(), nil, "Work", nil,
		domain_project.ProjectStatusActive, 1, nil, nil, now, now,
	)
}

func TestProject_Validate(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		mutate      func(p *domain_project.Project)
		expectedErr error
		checkResult func(t *testing.T, p *domain_project.Project)
	}{
		{
			name:        "valid project",
			mutate:      func(p *domain_project.Project) {},
			expectedErr: nil,
		},

		// --- Title ---
		{
			name:        "title empty is invalid",
			mutate:      func(p *domain_project.Project) { p.Title = "" },
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name:        "title exactly at lower boundary (3 chars) is valid",
			mutate:      func(p *domain_project.Project) { p.Title = "Hey" },
			expectedErr: nil,
		},
		{
			name:        "title one below lower boundary (2 chars) is invalid",
			mutate:      func(p *domain_project.Project) { p.Title = "He" },
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "title exactly at upper boundary (100 chars) is valid",
			mutate: func(p *domain_project.Project) {
				p.Title = strings.Repeat(
					"a", 100,
				)
			},
			expectedErr: nil,
		},
		{
			name: "title one above upper boundary (101 chars) is invalid",
			mutate: func(p *domain_project.Project) {
				p.Title = strings.Repeat(
					"a", 101,
				)
			},
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name:        "title only spaces (fails trim check)",
			mutate:      func(p *domain_project.Project) { p.Title = "     " },
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "title with leading/trailing spaces is trimmed before length check",
			mutate: func(p *domain_project.Project) {
				p.Title = "  Hi  " // trims to "Hi" -> 2 runes -> invalid
			},
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "title with multi-byte runes counted by rune, not by byte length",
			mutate: func(p *domain_project.Project) {
				p.Title = "Дом" // 3 runes, 6 bytes -> valid
			},
			expectedErr: nil,
		},
		{
			name: "title with multi-byte runes one above upper boundary (101 runes) is invalid",
			mutate: func(p *domain_project.Project) {
				p.Title = strings.Repeat("д", 101)
			},
			expectedErr: core_errors.ErrInvalidArgument,
		},

		// --- Position ---
		{
			name:        "position zero is invalid",
			mutate:      func(p *domain_project.Project) { p.Position = 0 },
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name:        "position negative is invalid",
			mutate:      func(p *domain_project.Project) { p.Position = -1 },
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name:        "position exactly at lower boundary (1) is valid",
			mutate:      func(p *domain_project.Project) { p.Position = 1 },
			expectedErr: nil,
		},

		// --- Timestamps ---
		{
			name:        "created_at is zero",
			mutate:      func(p *domain_project.Project) { p.CreatedAt = time.Time{} },
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name:        "updated_at is zero",
			mutate:      func(p *domain_project.Project) { p.UpdatedAt = time.Time{} },
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "updated_at before created_at",
			mutate: func(p *domain_project.Project) {
				p.UpdatedAt = p.CreatedAt.Add(-time.Hour)
			},
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "updated_at equal to created_at is valid",
			mutate: func(p *domain_project.Project) {
				p.UpdatedAt = p.CreatedAt
			},
			expectedErr: nil,
		},

		// --- Status ---
		{
			name:        "unknown status string is invalid",
			mutate:      func(p *domain_project.Project) { p.Status = domain_project.ProjectStatus("archived") },
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name:        "empty status string is invalid",
			mutate:      func(p *domain_project.Project) { p.Status = domain_project.ProjectStatus("") },
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name:        "status active with no completed_at is valid",
			mutate:      func(p *domain_project.Project) { p.Status = domain_project.ProjectStatusActive },
			expectedErr: nil,
		},
		{
			name:        "status dropped with no completed_at is valid",
			mutate:      func(p *domain_project.Project) { p.Status = domain_project.ProjectStatusDropped },
			expectedErr: nil,
		},
		{
			name: "status completed with completed_at set is valid",
			mutate: func(p *domain_project.Project) {
				p.Status = domain_project.ProjectStatusCompleted
				p.CompletedAt = ptr(p.CreatedAt.Add(time.Hour))
			},
			expectedErr: nil,
		},

		// --- Status / CompletedAt invariants ---
		{
			name: "status completed but completed_at nil is invalid",
			mutate: func(p *domain_project.Project) {
				p.Status = domain_project.ProjectStatusCompleted
				p.CompletedAt = nil
			},
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "status active but completed_at set is invalid",
			mutate: func(p *domain_project.Project) {
				p.Status = domain_project.ProjectStatusActive
				p.CompletedAt = ptr(p.CreatedAt.Add(time.Hour))
			},
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "status dropped but completed_at set is invalid",
			mutate: func(p *domain_project.Project) {
				p.Status = domain_project.ProjectStatusDropped
				p.CompletedAt = ptr(p.CreatedAt.Add(time.Hour))
			},
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "completed_at before created_at is invalid",
			mutate: func(p *domain_project.Project) {
				p.Status = domain_project.ProjectStatusCompleted
				p.CompletedAt = ptr(p.CreatedAt.Add(-time.Hour))
			},
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "completed_at equal to created_at is valid",
			mutate: func(p *domain_project.Project) {
				p.Status = domain_project.ProjectStatusCompleted
				p.CompletedAt = ptr(p.CreatedAt)
			},
			expectedErr: nil,
		},

		// --- AreaID ---
		{
			name:        "area_id nil is valid",
			mutate:      func(p *domain_project.Project) { p.AreaID = nil },
			expectedErr: nil,
		},
		{
			name:        "area_id valid non-nil UUID is valid",
			mutate:      func(p *domain_project.Project) { p.AreaID = ptr(uuid.New()) },
			expectedErr: nil,
		},
		{
			name:        "area_id set to empty (nil) UUID is invalid",
			mutate:      func(p *domain_project.Project) { p.AreaID = ptr(uuid.Nil) },
			expectedErr: core_errors.ErrInvalidArgument,
		},

		// --- Notes ---
		{
			name:        "notes nil is valid",
			mutate:      func(p *domain_project.Project) { p.Notes = nil },
			expectedErr: nil,
		},
		{
			name: "notes exactly at boundary (2000 chars) is valid",
			mutate: func(p *domain_project.Project) {
				p.Notes = ptr(
					strings.Repeat(
						"a", 2000,
					),
				)
			},
			expectedErr: nil,
		},
		{
			name: "notes one above boundary (2001 chars) is invalid",
			mutate: func(p *domain_project.Project) {
				p.Notes = ptr(
					strings.Repeat(
						"a", 2001,
					),
				)
			},
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "notes with multi-byte runes counted by rune, not by byte length",
			mutate: func(p *domain_project.Project) {
				p.Notes = ptr(
					strings.Repeat(
						"д", 2000,
					),
				) // 2000 runes, 4000 bytes -> valid
			},
			expectedErr: nil,
		},
		{
			name: "notes trimmed before length check",
			mutate: func(p *domain_project.Project) {
				padded := "  " + strings.Repeat("a", 2000) + "  "
				p.Notes = ptr(padded)
			},
			expectedErr: nil,
		},
		{
			name:        "notes only whitespace is valid (no minimum length)",
			mutate:      func(p *domain_project.Project) { p.Notes = ptr("     ") },
			expectedErr: nil,
		},

		// --- Deadline ---
		{
			name:        "deadline nil is valid",
			mutate:      func(p *domain_project.Project) { p.Deadline = nil },
			expectedErr: nil,
		},
		{
			name:        "deadline non-zero time is valid",
			mutate:      func(p *domain_project.Project) { p.Deadline = ptr(now.Add(24 * time.Hour)) },
			expectedErr: nil,
		},
		{
			name:        "deadline zero time (non-nil pointer) is invalid",
			mutate:      func(p *domain_project.Project) { p.Deadline = &time.Time{} },
			expectedErr: core_errors.ErrInvalidArgument,
		},

		// --- Version ---
		{
			name:   "negative version gets mutated to 0",
			mutate: func(p *domain_project.Project) { p.Version = -3 },
			checkResult: func(t *testing.T, p *domain_project.Project) {
				assert.Equal(t, 0, p.Version)
			},
			expectedErr: nil,
		},
		{
			name:   "zero version is left untouched",
			mutate: func(p *domain_project.Project) { p.Version = 0 },
			checkResult: func(t *testing.T, p *domain_project.Project) {
				assert.Equal(t, 0, p.Version)
			},
			expectedErr: nil,
		},
		{
			name:   "positive version is left untouched",
			mutate: func(p *domain_project.Project) { p.Version = 9 },
			checkResult: func(t *testing.T, p *domain_project.Project) {
				assert.Equal(t, 9, p.Version)
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				p := validProject(now)
				tt.mutate(&p)

				err := p.Validate()

				if tt.expectedErr != nil {
					require.Error(t, err)
					assert.ErrorIs(t, err, tt.expectedErr)
				} else {
					assert.NoError(t, err)
				}

				if tt.checkResult != nil {
					tt.checkResult(t, &p)
				}
			},
		)
	}
}

func TestNewProjectPatch(t *testing.T) {
	areaID := uuid.New()
	notes := "notes"
	deadline := time.Now().Add(24 * time.Hour)
	completedAt := time.Now()

	patch := domain_project.NewProjectPatch(
		nullableSet(&areaID),
		nullableSet("Title"),
		nullableSet(&notes),
		nullableSet(domain_project.ProjectStatusCompleted),
		nullableSet(&deadline),
		nullableSet(&completedAt),
	)

	assert.Equal(t, &areaID, *patch.AreaID.Value)
	assert.Equal(t, "Title", *patch.Title.Value)
	assert.Equal(t, &notes, *patch.Notes.Value)
	assert.Equal(t, domain_project.ProjectStatusCompleted, *patch.Status.Value)
	assert.Equal(t, &deadline, *patch.Deadline.Value)
	assert.Equal(t, &completedAt, *patch.CompletedAt.Value)
}

func TestProjectPatch_Validate(t *testing.T) {
	tests := []struct {
		name        string
		patch       domain_project.ProjectPatch
		expectedErr error
	}{
		{
			name:        "valid empty patch",
			patch:       domain_project.ProjectPatch{},
			expectedErr: nil,
		},
		{
			name:        "valid title patch",
			patch:       domain_project.ProjectPatch{Title: nullableSet("New Title")},
			expectedErr: nil,
		},
		{
			name:        "title Set true Value nil is invalid",
			patch:       domain_project.ProjectPatch{Title: nullableClear[string]()},
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name:        "status Set true Value nil is invalid",
			patch:       domain_project.ProjectPatch{Status: nullableClear[domain_project.ProjectStatus]()},
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name:        "status Set true with invalid value is invalid",
			patch:       domain_project.ProjectPatch{Status: nullableSet(domain_project.ProjectStatus("archived"))},
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name:        "status Set true with valid value is valid",
			patch:       domain_project.ProjectPatch{Status: nullableSet(domain_project.ProjectStatusDropped)},
			expectedErr: nil,
		},
		{
			name:        "area_id explicit clear (Value nil) is valid",
			patch:       domain_project.ProjectPatch{AreaID: nullableClear[*uuid.UUID]()},
			expectedErr: nil,
		},
		{
			name: "area_id explicit clear via nested nil pointer is valid",
			patch: domain_project.ProjectPatch{
				AreaID: domain.Nullable[*uuid.UUID]{
					Set: true, Value: ptr[*uuid.UUID](nil),
				},
			},
			expectedErr: nil,
		},
		{
			name:        "area_id set to a valid UUID is valid",
			patch:       domain_project.ProjectPatch{AreaID: nullableSet(ptr(uuid.New()))},
			expectedErr: nil,
		},
		{
			name:        "area_id set to empty (nil) UUID is invalid",
			patch:       domain_project.ProjectPatch{AreaID: nullableSet(ptr(uuid.Nil))},
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name:        "notes explicit clear (Value nil) is valid",
			patch:       domain_project.ProjectPatch{Notes: nullableClear[*string]()},
			expectedErr: nil,
		},
		{
			name:        "notes set within limit is valid",
			patch:       domain_project.ProjectPatch{Notes: nullableSet(ptr("some notes"))},
			expectedErr: nil,
		},
		{
			name: "notes set exactly at boundary (2000 chars) is valid",
			patch: domain_project.ProjectPatch{
				Notes: nullableSet(
					ptr(
						strings.Repeat(
							"a", 2000,
						),
					),
				),
			},
			expectedErr: nil,
		},
		{
			name: "notes set one above boundary (2001 chars) is invalid",
			patch: domain_project.ProjectPatch{
				Notes: nullableSet(
					ptr(
						strings.Repeat(
							"a", 2001,
						),
					),
				),
			},
			expectedErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "notes trimmed before length check",
			patch: domain_project.ProjectPatch{
				Notes: nullableSet(
					ptr(
						"  " + strings.Repeat(
							"a", 2000,
						) + "  ",
					),
				),
			},
			expectedErr: nil,
		},
		{
			name:        "deadline explicit clear (Value nil) is valid",
			patch:       domain_project.ProjectPatch{Deadline: nullableClear[*time.Time]()},
			expectedErr: nil,
		},
		{
			name:        "deadline set to a non-zero time is valid",
			patch:       domain_project.ProjectPatch{Deadline: nullableSet(ptr(time.Now().Add(24 * time.Hour)))},
			expectedErr: nil,
		},
		{
			name:        "deadline set to zero time is invalid",
			patch:       domain_project.ProjectPatch{Deadline: nullableSet(&time.Time{})},
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

// assertUnchanged verifies that every field of got matches want, except the
// ones explicitly excluded via the skip* flags. Used to prove ApplyPatch
// either applies a change atomically or leaves the project fully untouched.
func assertUnchangedExcept(
	t *testing.T,
	want, got domain_project.Project,
	skipTitle, skipStatus, skipAreaID, skipNotes, skipDeadline, skipCompletedAt bool,
) {
	t.Helper()

	if !skipTitle {
		assert.Equal(t, want.Title, got.Title, "Title should be unchanged")
	}
	if !skipStatus {
		assert.Equal(t, want.Status, got.Status, "Status should be unchanged")
	}
	if !skipAreaID {
		assert.Equal(t, want.AreaID, got.AreaID, "AreaID should be unchanged")
	}
	if !skipNotes {
		assert.Equal(t, want.Notes, got.Notes, "Notes should be unchanged")
	}
	if !skipDeadline {
		assert.Equal(
			t, want.Deadline, got.Deadline, "Deadline should be unchanged",
		)
	}
	if !skipCompletedAt {
		assert.Equal(
			t, want.CompletedAt, got.CompletedAt,
			"CompletedAt should be unchanged",
		)
	}

	assert.Equal(t, want.ID, got.ID)
	assert.Equal(t, want.UserID, got.UserID)
	assert.Equal(t, want.Position, got.Position)
	assert.Equal(t, want.CreatedAt, got.CreatedAt)
	assert.Equal(t, want.UpdatedAt, got.UpdatedAt)
}

func baseActiveProject(now time.Time) domain_project.Project {
	return domain_project.NewProject(
		uuid.New(), 1, uuid.New(), nil, "Initial Title", nil,
		domain_project.ProjectStatusActive, 1, nil, nil, now, now,
	)
}

func TestProject_ApplyPatch(t *testing.T) {
	now := time.Now()

	t.Run(
		"successfully apply title patch", func(t *testing.T) {
			original := baseActiveProject(now)
			p := original
			err := p.ApplyPatch(domain_project.ProjectPatch{Title: nullableSet("Updated Title")})

			require.NoError(t, err)
			assert.Equal(t, "Updated Title", p.Title)
			assertUnchangedExcept(
				t, original, p, true, false, false, false, false, false,
			)
		},
	)

	t.Run(
		"empty patch changes nothing", func(t *testing.T) {
			original := baseActiveProject(now)
			p := original
			err := p.ApplyPatch(domain_project.ProjectPatch{})

			require.NoError(t, err)
			assert.Equal(t, original, p)
		},
	)

	t.Run(
		"fail: patch itself invalid (title Set true, Value nil) — project fully unchanged",
		func(t *testing.T) {
			original := baseActiveProject(now)
			p := original
			err := p.ApplyPatch(domain_project.ProjectPatch{Title: nullableClear[string]()})

			require.Error(t, err)
			assert.ErrorIs(t, err, core_errors.ErrInvalidArgument)
			assert.Equal(t, original, p)
		},
	)

	t.Run(
		"fail: patched title too short — project fully unchanged",
		func(t *testing.T) {
			original := baseActiveProject(now)
			p := original
			err := p.ApplyPatch(domain_project.ProjectPatch{Title: nullableSet("No")})

			require.Error(t, err)
			assert.ErrorIs(t, err, core_errors.ErrInvalidArgument)
			assert.Equal(t, original, p)
		},
	)

	t.Run(
		"area_id: set from nil to a concrete area", func(t *testing.T) {
			original := baseActiveProject(now)
			areaID := uuid.New()
			p := original
			err := p.ApplyPatch(domain_project.ProjectPatch{AreaID: nullableSet(ptr(areaID))})

			require.NoError(t, err)
			require.NotNil(t, p.AreaID)
			assert.Equal(t, areaID, *p.AreaID)
		},
	)

	t.Run(
		"area_id: explicit clear (Value nil) removes area", func(t *testing.T) {
			areaID := uuid.New()
			original := baseActiveProject(now)
			original.AreaID = &areaID
			p := original
			err := p.ApplyPatch(domain_project.ProjectPatch{AreaID: nullableClear[*uuid.UUID]()})

			require.NoError(t, err)
			assert.Nil(t, p.AreaID)
		},
	)

	t.Run(
		"area_id: explicit clear via nested nil pointer removes area",
		func(t *testing.T) {
			areaID := uuid.New()
			original := baseActiveProject(now)
			original.AreaID = &areaID
			p := original
			err := p.ApplyPatch(
				domain_project.ProjectPatch{
					AreaID: domain.Nullable[*uuid.UUID]{
						Set: true, Value: ptr[*uuid.UUID](nil),
					},
				},
			)

			require.NoError(t, err)
			assert.Nil(t, p.AreaID)
		},
	)

	t.Run(
		"area_id: patching to empty UUID fails before mutation",
		func(t *testing.T) {
			original := baseActiveProject(now)
			p := original
			err := p.ApplyPatch(domain_project.ProjectPatch{AreaID: nullableSet(ptr(uuid.Nil))})

			require.Error(t, err)
			assert.ErrorIs(t, err, core_errors.ErrInvalidArgument)
			assert.Equal(t, original, p)
		},
	)

	t.Run(
		"notes: set, then explicit clear", func(t *testing.T) {
			original := baseActiveProject(now)
			p := original
			err := p.ApplyPatch(domain_project.ProjectPatch{Notes: nullableSet(ptr("some notes"))})
			require.NoError(t, err)
			require.NotNil(t, p.Notes)
			assert.Equal(t, "some notes", *p.Notes)

			err = p.ApplyPatch(domain_project.ProjectPatch{Notes: nullableClear[*string]()})
			require.NoError(t, err)
			assert.Nil(t, p.Notes)
		},
	)

	t.Run(
		"notes: patching over the limit fails before mutation",
		func(t *testing.T) {
			original := baseActiveProject(now)
			p := original
			err := p.ApplyPatch(
				domain_project.ProjectPatch{
					Notes: nullableSet(
						ptr(
							strings.Repeat(
								"a", 2001,
							),
						),
					),
				},
			)

			require.Error(t, err)
			assert.ErrorIs(t, err, core_errors.ErrInvalidArgument)
			assert.Equal(t, original, p)
		},
	)

	t.Run(
		"deadline: set, then explicit clear", func(t *testing.T) {
			original := baseActiveProject(now)
			deadline := now.Add(48 * time.Hour)
			p := original
			err := p.ApplyPatch(domain_project.ProjectPatch{Deadline: nullableSet(&deadline)})
			require.NoError(t, err)
			require.NotNil(t, p.Deadline)
			assert.Equal(t, deadline, *p.Deadline)

			err = p.ApplyPatch(domain_project.ProjectPatch{Deadline: nullableClear[*time.Time]()})
			require.NoError(t, err)
			assert.Nil(t, p.Deadline)
		},
	)

	t.Run(
		"deadline: patching to zero time fails before mutation",
		func(t *testing.T) {
			original := baseActiveProject(now)
			p := original
			err := p.ApplyPatch(domain_project.ProjectPatch{Deadline: nullableSet(&time.Time{})})

			require.Error(t, err)
			assert.ErrorIs(t, err, core_errors.ErrInvalidArgument)
			assert.Equal(t, original, p)
		},
	)

	// --- Status / CompletedAt interaction — the trickiest part of ApplyPatch ---

	t.Run(
		"status -> completed without explicit completed_at auto-sets completed_at to now",
		func(t *testing.T) {
			original := baseActiveProject(now)
			p := original
			before := time.Now().UTC()
			err := p.ApplyPatch(domain_project.ProjectPatch{Status: nullableSet(domain_project.ProjectStatusCompleted)})
			after := time.Now().UTC()

			require.NoError(t, err)
			assert.Equal(t, domain_project.ProjectStatusCompleted, p.Status)
			require.NotNil(t, p.CompletedAt)
			assert.False(t, p.CompletedAt.Before(before))
			assert.False(t, p.CompletedAt.After(after))
		},
	)

	t.Run(
		"status -> completed with explicit completed_at uses the explicit value",
		func(t *testing.T) {
			original := baseActiveProject(now)
			explicitCompletedAt := now.Add(time.Hour)
			p := original
			err := p.ApplyPatch(
				domain_project.ProjectPatch{
					Status:      nullableSet(domain_project.ProjectStatusCompleted),
					CompletedAt: nullableSet(&explicitCompletedAt),
				},
			)

			require.NoError(t, err)
			require.NotNil(t, p.CompletedAt)
			assert.Equal(t, explicitCompletedAt, *p.CompletedAt)
		},
	)

	t.Run(
		"status -> completed with explicit completed_at cleared fails validation",
		func(t *testing.T) {
			original := baseActiveProject(now)
			p := original
			err := p.ApplyPatch(
				domain_project.ProjectPatch{
					Status:      nullableSet(domain_project.ProjectStatusCompleted),
					CompletedAt: nullableClear[*time.Time](),
				},
			)

			require.Error(t, err)
			assert.ErrorIs(t, err, core_errors.ErrInvalidArgument)
			assert.Equal(t, original, p)
		},
	)

	t.Run(
		"status: completed -> active without explicit completed_at auto-clears completed_at",
		func(t *testing.T) {
			completedAt := now.Add(time.Hour)
			original := baseActiveProject(now)
			original.Status = domain_project.ProjectStatusCompleted
			original.CompletedAt = &completedAt
			p := original

			err := p.ApplyPatch(domain_project.ProjectPatch{Status: nullableSet(domain_project.ProjectStatusActive)})

			require.NoError(t, err)
			assert.Equal(t, domain_project.ProjectStatusActive, p.Status)
			assert.Nil(t, p.CompletedAt)
		},
	)

	t.Run(
		"status: completed -> dropped without explicit completed_at auto-clears completed_at",
		func(t *testing.T) {
			completedAt := now.Add(time.Hour)
			original := baseActiveProject(now)
			original.Status = domain_project.ProjectStatusCompleted
			original.CompletedAt = &completedAt
			p := original

			err := p.ApplyPatch(domain_project.ProjectPatch{Status: nullableSet(domain_project.ProjectStatusDropped)})

			require.NoError(t, err)
			assert.Equal(t, domain_project.ProjectStatusDropped, p.Status)
			assert.Nil(t, p.CompletedAt)
		},
	)

	t.Run(
		"setting completed_at without changing status away from active fails validation",
		func(t *testing.T) {
			original := baseActiveProject(now) // status: active
			p := original
			explicitCompletedAt := now.Add(time.Hour)
			err := p.ApplyPatch(domain_project.ProjectPatch{CompletedAt: nullableSet(&explicitCompletedAt)})

			require.Error(t, err)
			assert.ErrorIs(t, err, core_errors.ErrInvalidArgument)
			assert.Equal(t, original, p)
		},
	)

	t.Run(
		"re-patching status to the same completed value with explicit completed_at keeps explicit value",
		func(t *testing.T) {
			completedAt := now.Add(time.Hour)
			original := baseActiveProject(now)
			original.Status = domain_project.ProjectStatusCompleted
			original.CompletedAt = &completedAt
			p := original

			newCompletedAt := now.Add(2 * time.Hour)
			err := p.ApplyPatch(
				domain_project.ProjectPatch{
					Status:      nullableSet(domain_project.ProjectStatusCompleted),
					CompletedAt: nullableSet(&newCompletedAt),
				},
			)

			require.NoError(t, err)
			require.NotNil(t, p.CompletedAt)
			assert.Equal(t, newCompletedAt, *p.CompletedAt)
		},
	)

	t.Run(
		"patching completed_at to zero time fails validation (caught via Before(created_at), not patch.Validate)",
		func(t *testing.T) {
			original := baseActiveProject(now)
			original.Status = domain_project.ProjectStatusCompleted
			original.CompletedAt = ptr(now.Add(time.Hour))
			p := original

			err := p.ApplyPatch(
				domain_project.ProjectPatch{
					CompletedAt: domain.Nullable[*time.Time]{
						Set: true, Value: ptr(&time.Time{}),
					},
				},
			)

			require.Error(t, err)
			assert.ErrorIs(t, err, core_errors.ErrInvalidArgument)
			assert.Equal(t, original, p)
		},
	)

	t.Run(
		"changing only status leaves title/notes/area/deadline untouched",
		func(t *testing.T) {
			areaID := uuid.New()
			notes := "keep me"
			deadline := now.Add(72 * time.Hour)
			original := baseActiveProject(now)
			original.AreaID = &areaID
			original.Notes = &notes
			original.Deadline = &deadline
			p := original

			err := p.ApplyPatch(domain_project.ProjectPatch{Status: nullableSet(domain_project.ProjectStatusDropped)})

			require.NoError(t, err)
			assert.Equal(t, domain_project.ProjectStatusDropped, p.Status)
			assertUnchangedExcept(
				t, original, p, false, true, false, false, false, true,
			)
		},
	)
}
