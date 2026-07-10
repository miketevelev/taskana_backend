package core_recurrence

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
	"github.com/teambition/rrule-go"
)

func NextFixedDate(rule string, after time.Time) (time.Time, error) {
	r, err := rrule.StrToRRule(rule)
	if err != nil {
		return time.Time{}, fmt.Errorf(
			"parse rrule: %w", core_errors.ErrInvalidArgument,
		)
	}

	next := r.After(after, false)
	if next.IsZero() {
		return time.Time{}, fmt.Errorf(
			"no next occurrence: %w", core_errors.ErrInvalidArgument,
		)
	}

	return truncateDate(next), nil
}

func NextFromCompletionDate(rule string, completedAt time.Time) (
	time.Time,
	error,
) {
	days, err := parseIntervalDays(rule)
	if err != nil {
		return time.Time{}, err
	}

	next := truncateDate(completedAt).AddDate(0, 0, days)
	return next, nil
}

func parseIntervalDays(rule string) (int, error) {
	rule = strings.TrimSpace(rule)
	if strings.HasPrefix(rule, "INTERVAL:") {
		days, err := strconv.Atoi(strings.TrimPrefix(rule, "INTERVAL:"))
		if err != nil || days <= 0 {
			return 0, fmt.Errorf(
				"invalid interval rule: %w", core_errors.ErrInvalidArgument,
			)
		}
		return days, nil
	}

	var days int
	if _, err := fmt.Sscanf(rule, "P%dD", &days); err == nil && days > 0 {
		return days, nil
	}

	return 0, fmt.Errorf(
		"unsupported from_completion rule: %w", core_errors.ErrInvalidArgument,
	)
}

func truncateDate(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}
