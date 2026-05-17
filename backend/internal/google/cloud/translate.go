package cloud

import (
	"context"
	"fmt"

	translate "cloud.google.com/go/translate/apiv3"
	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/translate/apiv3/translatepb"
	"github.com/koo-arch/adjusta-backend/configs"
)

type serviceAccount struct {
	ProjectID string `json:"project_id"`
}

type Translator struct {
	client *translate.TranslationClient
	projectID string
}

func NewTranslator(ctx context.Context) (*Translator, error) {
	client, err := translate.NewTranslationClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}

	var projectID string
	if configs.GetEnv("GO_ENV") == "development" {
		projectID, err = GetProjectID()
		if err != nil {
			return nil, fmt.Errorf("failed to get project ID: %v", err)
		}
	} else {
		projectID, err = metadata.ProjectIDWithContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get project ID: %v", err)
		}
	}

	return &Translator{
		client:    client,
		projectID: projectID,
	}, nil
}

func (t *Translator) Close() error {
	return t.client.Close()
}

func (t *Translator) getParent() string {
	return fmt.Sprintf("projects/%s/locations/global", t.projectID)
}

func (t *Translator) TranslateText(ctx context.Context, text string, targetLang string) (string, error) {

	request := &translatepb.TranslateTextRequest{
		Parent: t.getParent(),
		Contents: []string{text},
		TargetLanguageCode: targetLang,
	}

	resp, err := t.client.TranslateText(ctx, request)
	if err != nil {
		return "", fmt.Errorf("failed to translate text: %v", err)
	}

	if len(resp.Translations) == 0 {
		return "", fmt.Errorf("no translations returned")
	}
	
	return resp.GetTranslations()[0].GetTranslatedText(), nil
}
