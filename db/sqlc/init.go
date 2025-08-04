package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DBPool *pgxpool.Pool

func Connect() {
	connStr := "postgres://postgres:root@localhost:5432/kuhakuanime?sslmode=disable"

	var err error
	DBPool, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Connected to database")
}