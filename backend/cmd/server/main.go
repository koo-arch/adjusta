package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/koo-arch/adjusta-backend/internal/app"
	"github.com/koo-arch/adjusta-backend/internal/config"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx, config.New()); err != nil {
		log.Fatal(err)
	}
}
