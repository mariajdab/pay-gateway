package server

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	server *http.Server
}

func NewApp() *App {
	dbStr, connected := os.LookupEnv("DB_SOURCE")
	if !connected {
		log.Fatalf("Failed to read database connection")
	}

	_, err := pgxpool.New(context.Background(), dbStr)
	if err != nil {
		log.Fatalf("Unable to connect to db: %v", err)
	}

	return &App{}
}
