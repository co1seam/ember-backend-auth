package main

import (
	"context"
	"flag"
	"github.com/co1seam/ember-backend-auth/config"
	"github.com/co1seam/ember-backend-auth/internal/adapters/repository"
	"github.com/co1seam/ember-backend-auth/internal/adapters/rpc"
	"github.com/co1seam/ember-backend-auth/internal/core/models"
	"github.com/co1seam/ember-backend-auth/internal/core/services"
	"github.com/co1seam/ember-backend-auth/pkg/logger"
	"log"
	"log/slog"
	"os"
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

	cache := repository.NewRedis(cfg.Redis.Host, cfg.Redis.Port)

	opts := &models.Options{
		Logger: log,
		Config: cfg,
	}

	repos := repository.NewRepository(db.DB, cache, opts)
	service := services.NewService(repos, opts)
	handler := rpc.NewHandler(service, opts)

	server := rpc.NewServer()
	if err := server.Run(handler); err != nil {
		return
	}

	if err := db.DB.Close(); err != nil {
		log.Error("error: ", err)
	}
}
