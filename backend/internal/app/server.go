package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/koo-arch/adjusta-backend/internal/config"
	infraDatabase "github.com/koo-arch/adjusta-backend/internal/infrastructure/database"
)

const shutdownTimeout = 10 * time.Second

func Run(ctx context.Context, cfg config.Config) (runErr error) {
	client, err := infraDatabase.New(cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer func() {
		if err := client.Close(); err != nil && runErr == nil {
			runErr = fmt.Errorf("failed closing connection to postgres: %w", err)
		}
	}()

	if cfg.AutoMigrate {
		if err := infraDatabase.Migrate(ctx, client); err != nil {
			return err
		}
	}

	deps := buildDependencies(client, cfg)
	router := newRouter(cfg, deps)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	errCh := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("failed to run server: %w", err)
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("failed to shutdown server: %w", err)
		}
		return nil
	}
}
