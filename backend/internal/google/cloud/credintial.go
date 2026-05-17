package cloud

import (
	"os"
	"fmt"
	"log"
	"encoding/json"

	"github.com/koo-arch/adjusta-backend/configs"
)

func GetProjectID() (string, error) {
	credsPath := configs.GetEnv("GOOGLE_APPLICATION_CREDENTIALS")
	if credsPath == "" {
		return "", fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS is not set")
	}

	file, err := os.Open(credsPath)
	if err != nil {
		return "", fmt.Errorf("failed to open service account file: %v", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("failed to close service account file: %v", err)
		}
	}()

	var sa serviceAccount
	if err := json.NewDecoder(file).Decode(&sa); err != nil {
		return "", fmt.Errorf("failed to decode service account file: %v", err)
	}

	return sa.ProjectID, nil
}