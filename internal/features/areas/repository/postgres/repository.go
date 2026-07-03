package areas_postgres_repository

import (
	core_postgres_pool "github.com/miketevelev/taskana_backend/internal/core/repository/postgres/pool"
)

type AreasRepository struct {
	pool core_postgres_pool.Pool
}

func NewAreasRepository(
	pool core_postgres_pool.Pool,
) *AreasRepository {
	return &AreasRepository{
		pool: pool,
	}
}
