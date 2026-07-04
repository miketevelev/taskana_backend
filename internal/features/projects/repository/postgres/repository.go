package projects_postgres_repository

import core_postgres_pool "github.com/miketevelev/taskana_backend/internal/core/repository/postgres/pool"

type ProjectRepository struct {
	pool core_postgres_pool.Pool
}

func NewProjectRepository(
	pool core_postgres_pool.Pool,
) *ProjectRepository {
	return &ProjectRepository{
		pool: pool,
	}
}
