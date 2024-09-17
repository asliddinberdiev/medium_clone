package main

import (
	L "log"

	"github.com/asliddinberdiev/medium_clone/config"
	"github.com/asliddinberdiev/medium_clone/handler"
	"github.com/asliddinberdiev/medium_clone/repository"
	"github.com/asliddinberdiev/medium_clone/server"
	"github.com/asliddinberdiev/medium_clone/service"
	"go.uber.org/zap"

	_ "github.com/lib/pq"
)

func main() {
	log := service.InitLogger("logs", "app")
	defer func() {
		if err := log.Sync(); err != nil {
			L.Fatalf("logger sync error: %v\n", err)
		}
	}()

	cfg := config.Load(".", log)

	// initialize db
	db, err := repository.NewPostgresDB(repository.PostgresConfig{
		User:     cfg.Postgres.User,
		Password: cfg.Postgres.Password,
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		Database: cfg.Postgres.Database,
		SSLMode:  cfg.Postgres.SSLMode,
	}, log)
	if err != nil {
		log.Fatal("failed to initialize db", zap.Error(err))
	}

	// initialize repository
	repos := repository.NewRepository(db, log)

	// initialize services
	services := service.NewService(repos, log)

	// initialize handlers
	handlers := handler.NewHandler(services, cfg.App.Version, log)

	log.Info("app run", zap.String("port", cfg.App.Port))
	srv := new(server.Server)
	if err := srv.Run(cfg.App.Port, handlers.InitRoutes()); err != nil {
		log.Fatal("error occurred while running http server", zap.Error(err))
	}
}
