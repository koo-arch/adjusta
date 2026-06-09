package utils

import (
	"context"
	"fmt"

	"github.com/gosimple/slug"
)

func NormalizeToSlug(ctx context.Context, text string) (string, error) {
	translated, err := TranslateText(ctx, text, "en")
	if err != nil {
		return "", fmt.Errorf("failed to translate text: %v", err)
	}
	text = translated
	fmt.Println("text: ", text)
	return slug.Make(text), nil
}

func EnsureUniqueSlug(ctx context.Context, existingSlugs map[string]struct{}, baseSlug string, strLength int) string {
	if _, exists := existingSlugs[baseSlug]; !exists {
		return baseSlug
	}
	for {
		newSlug := fmt.Sprintf("%s-%s", baseSlug, GenerateRamdomString(strLength))
		if _, exists := existingSlugs[newSlug]; !exists {
			return newSlug
		}
	}
}
