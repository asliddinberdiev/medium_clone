package main

import (
	"log"
	"os"
	"path"

	"github.com/asliddinberdiev/medium_clone/config"
	"github.com/asliddinberdiev/medium_clone/handler"
	"github.com/asliddinberdiev/medium_clone/repository"
	"github.com/asliddinberdiev/medium_clone/server"
	"github.com/asliddinberdiev/medium_clone/service"

	_ "github.com/lib/pq"
)

func main() {
	_ = os.Mkdir("logs", 0770)
	logFile, err := os.OpenFile(path.Join("logs", "app.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("create log file error: ", err)
	}
	log.SetOutput(logFile)

	cfg := config.Load(".")

	db, err := repository.NewPostgresDB(repository.PostgresConfig{
		User:     cfg.Postgres.User,
		Password: cfg.Postgres.Password,
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		Database: cfg.Postgres.Database,
		SSLMode:  cfg.Postgres.SSLMode,
	})
	if err != nil {
		log.Fatalln("failed to initialize db error: ", err)
	}

	// initialize rdb
	rdb, err := repository.NewRedisDB(repository.RedisConfig{
		Host:     cfg.Redis.Host,
		Password: cfg.Redis.Password,
		Port:     cfg.Redis.Port,
	})
	if err != nil {
		log.Fatalln("failed to initialize rsdb error: ", err)
	}

	repos := repository.NewRepository(db, rdb)
	services := service.NewService(repos, cfg.App)
	handlers := handler.NewHandler(services, cfg.App)

	log.Println("app run on port: ", cfg.App.Port)
	srv := new(server.Server)
	if err := srv.Run(cfg.App.Port, handlers.InitRoutes()); err != nil {
		log.Fatalln("running http server error: ", err)
	}
}
