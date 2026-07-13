package timetracking_postgres_repository

import core_postgres_pool "github.com/miketevelev/taskana_backend/internal/core/repository/postgres/pool"

type TimeTrackingRepository struct {
	pool core_postgres_pool.Pool
}

func NewTimeTrackingRepository(
	pool core_postgres_pool.Pool,
) *TimeTrackingRepository {
	return &TimeTrackingRepository{
		pool: pool,
	}
}
