package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/koo-arch/adjusta-backend/internal/app"
	"github.com/koo-arch/adjusta-backend/internal/config"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.NewServer()
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Run(ctx, cfg); err != nil {
		log.Fatal(err)
	}
}
