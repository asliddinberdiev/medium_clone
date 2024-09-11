package main

import (
	"fmt"
	"log"
	"net/http"

	router "github.com/asliddinberdiev/medium_clone/api"
	"github.com/asliddinberdiev/medium_clone/config"
	"github.com/asliddinberdiev/medium_clone/storage"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load(".")

	psqlUrl := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User,
		cfg.Postgres.Password, cfg.Postgres.Database)

	psqlConn, err := sqlx.Connect("postgres", psqlUrl)
	if err != nil {
		log.Fatalf("Failed to connect to postgres: %v", err)
	}

	strg := storage.NewStorage(psqlConn)

	router := router.NewRouter(&router.Options{Strg: strg})

	server := &http.Server{
		Addr:    ":" + cfg.App.Port,
		Handler: router,
	}

	log.Printf("Start to run to server: %v\n", cfg.App.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to run to server: %v", err)
	}

}
