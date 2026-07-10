package checklists_postgres_repository

import core_postgres_pool "github.com/miketevelev/taskana_backend/internal/core/repository/postgres/pool"

type ChecklistsRepository struct {
	pool core_postgres_pool.Pool
}

func NewChecklistsRepository(
	pool core_postgres_pool.Pool,
) *ChecklistsRepository {
	return &ChecklistsRepository{
		pool: pool,
	}
}
