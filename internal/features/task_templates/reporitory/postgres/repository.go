package task_templates_postgres_repository

import core_postgres_pool "github.com/miketevelev/taskana_backend/internal/core/repository/postgres/pool"

type TaskTemplateRepository struct {
	pool core_postgres_pool.Pool
}

func NewTaskTemplateRepository(
	pool core_postgres_pool.Pool,
) *TaskTemplateRepository {
	return &TaskTemplateRepository{
		pool: pool,
	}
}
