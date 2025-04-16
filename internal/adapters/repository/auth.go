package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/co1seam/ember-backend-auth/internal/core/models"
)

type Authorization struct {
	db   *sql.DB
	opts *models.Options
}

func NewAuthorization(db *sql.DB, opts *models.Options) *Authorization {
	return &Authorization{
		db:   db,
		opts: opts,
	}
}

func (a *Authorization) Create(ctx context.Context, entity ...interface{}) (interface{}, error) {
	var id int
	user := entity[0].(models.SignUpRequest)

	query := fmt.Sprintf("INSERT INTO %s (user_name,user_email,user_password) VALUES ($1, $2, $3) RETURNING user_id", models.UserTable)
	err := a.db.QueryRowContext(ctx, query, user.Name, user.Email, user.Password).Scan(&id)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (a *Authorization) Read(ctx context.Context, entity ...interface{}) (interface{}, error) {
	var id int
	request := entity[0].(models.SignInRequest)

	query := fmt.Sprintf("SELECT user_id FROM %s WHERE user_email = $1 AND user_password = $2", models.UserTable)
	err := a.db.QueryRowContext(ctx, query, request.Email, request.Password).Scan(&id)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (a *Authorization) Update(ctx context.Context, entity ...interface{}) (interface{}, error) {
	return nil, nil
}

func (a *Authorization) Delete(ctx context.Context, filter ...interface{}) (interface{}, error) {
	return nil, nil
}
