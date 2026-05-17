package utils

import (
	"context"
	"fmt"
	"log"

	"github.com/koo-arch/adjusta-backend/internal/google/cloud"
)

func TranslateText(ctx context.Context, text string, targetLang string) (string, error) {
	t, err := cloud.NewTranslator(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create translator: %v", err)
	}

	defer func() {
		if err := t.Close(); err != nil {
			log.Printf("failed closing translator: %v", err)
		}
	}()

	translated, err := t.TranslateText(ctx, text, targetLang)
	if err != nil {
		return "", fmt.Errorf("failed to translate text: %v", err)
	}

	return translated, nil
}