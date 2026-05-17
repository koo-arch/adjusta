package utils

import (
	"fmt"
	"context"

	"github.com/gosimple/slug"
)


func NormalizeToSlug(ctx context.Context, summary string) (string, error) {
	// if !IsEnglish(summary) {
		translated, err := TranslateText(ctx, summary, "en")
		if err != nil {
			return "", fmt.Errorf("failed to translate summary: %v", err)
		}
		summary = translated
	// }
	fmt.Println("summary: ", summary)
	return slug.Make(summary), nil
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