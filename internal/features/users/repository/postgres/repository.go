package user_postgres_repository

import core_postgres_pool "github.com/miketevelev/taskana_backend/internal/core/repository/postgres/pool"

type UserRepository struct {
	pool core_postgres_pool.Pool
}

func NewUserRepository(
	pool core_postgres_pool.Pool,
) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}
