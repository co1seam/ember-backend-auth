package models

import (
	"github.com/co1seam/ember-backend-auth/internal/config"
	"github.com/co1seam/ember-backend-auth/pkg/logger"
)

type Options struct {
	Logger *logger.Logger
	Config *config.Config
}

const (
	UserTable = "users"
)
