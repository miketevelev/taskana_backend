package analytics_postgres_repository

import core_postgres_pool "github.com/miketevelev/taskana_backend/internal/core/repository/postgres/pool"

type AnalyticsRepository struct {
	pool core_postgres_pool.Pool
}

func NewAnalyticsRepository(
	pool core_postgres_pool.Pool,
) *AnalyticsRepository {
	return &AnalyticsRepository{
		pool: pool,
	}
}
