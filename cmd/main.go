package main

import (
	"context"
	"fmt"
	"log"

	"github.com/asliddinberdiev/medium_clone/config"
	"github.com/asliddinberdiev/medium_clone/storage"
	"github.com/asliddinberdiev/medium_clone/storage/repo"
	"github.com/google/uuid"
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

	id, err := uuid.NewRandom()
	if err != nil {
		log.Fatal(err)
	}

	user, err := strg.User().Create(context.TODO(), &repo.User{
		ID:        id.String(),
		FirstName: "Asliddin",
		LastName:  "Berdiev",
		Email:     "asliddinwork@gmail.com",
		Password:  "12345678",
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(user)
}
