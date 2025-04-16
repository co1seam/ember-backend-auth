package repository

import (
	"database/sql"
	"github.com/co1seam/ember-backend-auth/internal/core/models"
	"github.com/co1seam/ember-backend-auth/internal/ports"
)

type Repository struct {
	Authorization ports.IAuthRepo
	opts          *models.Options
}

func NewRepository(db *sql.DB, opts *models.Options) *Repository {
	return &Repository{
		Authorization: NewAuthorization(db, opts),
		opts:          opts,
	}
}
