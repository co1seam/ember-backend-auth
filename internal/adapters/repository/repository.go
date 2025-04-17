package repository

import (
	"database/sql"
	"github.com/co1seam/ember-backend-auth/internal/core/models"
	"github.com/co1seam/ember-backend-auth/internal/ports"
)

type Repository struct {
	Authorization ports.IAuthRepo
	Cache         *Redis
}

func NewRepository(db *sql.DB, cache *Redis, opts *models.Options) *Repository {
	return &Repository{
		Authorization: NewAuthorization(db, opts),
		Cache:         cache,
	}
}
