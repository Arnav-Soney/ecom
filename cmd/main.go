package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/Arnav-Soney/ecom/internal/env"
	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()
	cfg := Config{
		addr: ":8080",
		db: dbConfig{
			dsn: env.GetString("GOOSE_DBSTRING", "host= localhost user=postgres password=postgres dbname=ecom port=5432 sslmode=disable"),
		},
	}

	// logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Database
	conn, err := pgx.Connect(ctx, cfg.db.dsn)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)
	
	logger.Info("connected to databse", "dsn", cfg.db.dsn)

	api := application{
		config: cfg,
		db: conn, 
	}
	if err := api.run(api.mount()); err != nil {
		log.Printf("Server has failed to start, err: %s", err)
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}

}
