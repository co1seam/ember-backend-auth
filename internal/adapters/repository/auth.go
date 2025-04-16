package postgres

import (
	"context"
	"database/sql"
	"reflect"
)

type Authorization struct {
	db *sql.DB
}

func (a *Authorization) Create(ctx context.Context, request ...interface{}) error {
	email := reflect.ValueOf(request)
	_, err := a.db.QueryContext(ctx, "INSERT INTO user (user_email) VALUES ($1)", email)
	if err != nil {
		return err
	}
	return nil
}
