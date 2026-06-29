package core_postgres_pool

import (
	"errors"
)

var (
	ErrUnknown           = errors.New("unknown")
	ErrNoRows            = errors.New("no rows")
	ErrViolateForeignKey = errors.New(
		"violate foreign key",
	)
)
