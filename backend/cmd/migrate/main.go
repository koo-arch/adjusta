package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/koo-arch/adjusta-backend/internal/config"
	"github.com/koo-arch/adjusta-backend/internal/infrastructure/database"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}

	client, err := database.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	if err := database.Migrate(ctx, client); err != nil {
		log.Fatal(err)
	}

	log.Println("migration completed")
}
