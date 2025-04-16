package main

import (
	"context"
	"flag"
	"github.com/co1seam/ember-backend-auth/internal/adapters/http"
	"github.com/co1seam/ember-backend-auth/internal/adapters/repository"
	"github.com/co1seam/ember-backend-auth/internal/config"
	"github.com/co1seam/ember-backend-auth/internal/core/models"
	"github.com/co1seam/ember-backend-auth/internal/core/services"
	"github.com/co1seam/ember-backend-auth/pkg/logger"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()

	cfgFlag := flag.String("config", "", "flag to add config path")

	flag.Parse()

	cfg, err := config.New(cfgFlag)
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
		err := log.Error("error: ", err)
		if err != nil {
			return
		}
	}

	opts := &models.Options{
		Logger: log,
		Config: cfg,
	}

	repos := repository.NewRepository(db.DB, opts)
	service := services.NewService(repos, opts)
	handler := http.NewHandler(service, opts)

	server := http.NewServer()
	handler.Router(server.Server)
	if err := server.Run("8080"); err != nil {
		return
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := server.Shutdown(); err != nil {
		log.Error("error: ", err)
	}

	if err := db.DB.Close(); err != nil {
		log.Error("error: ", err)
	}
}
