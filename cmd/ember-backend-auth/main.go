package main

import (
	"context"
	"github.com/co1seam/ember-backend-auth/internal/adapters/http"
	"github.com/co1seam/ember-backend-auth/internal/adapters/repository"
	"github.com/co1seam/ember-backend-auth/internal/config"
	"github.com/co1seam/ember-backend-auth/internal/core/models"
	"github.com/co1seam/ember-backend-auth/internal/core/services"
	"github.com/co1seam/ember-backend-auth/pkg/logger"
	"log"
	"log/slog"
	"os"
)

func main() {
	ctx := context.Background()

	cfg, err := config.New("")
	if err != nil {
		log.Fatal(err)
	}

	log := logger.New(ctx, logger.Options{
		Level:     slog.LevelError,
		AddSource: true,
		JSON:      true,
		Output:    os.Stdout,
	})

	db, err := repository.NewPostgres(ctx, &cfg.Database)
	if err != nil {
		return
	}

	_ = &models.Options{
		Logger: log,
		Config: cfg,
	}

	repos := repository.NewRepository(db.DB)
	service := services.NewService(repos)
	handler := http.NewHandler(service)

	server := http.NewServer()
	handler.Router(server.Server)
	if err := server.Run("8080"); err != nil {
		return
	}
}
