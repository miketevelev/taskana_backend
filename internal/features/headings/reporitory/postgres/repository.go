package heading_postgres_repository

import core_postgres_pool "github.com/miketevelev/taskana_backend/internal/core/repository/postgres/pool"

type HeadingRepository struct {
	pool core_postgres_pool.Pool
}

func NewHeadingRepository(
	pool core_postgres_pool.Pool,
) *HeadingRepository {
	return &HeadingRepository{
		pool: pool,
	}
}
