package database

import (
	"fmt"

	"github.com/koo-arch/adjusta-backend/internal/infrastructure/ent"

	_ "github.com/koo-arch/adjusta-backend/internal/infrastructure/ent/runtime"
	_ "github.com/lib/pq"
)

func New(databaseURL string) (*ent.Client, error) {
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	client, err := ent.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to postgres: %w", err)
	}

	return client, nil
}
