package googleapierror

import (
	stderrors "errors"

	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"google.golang.org/api/googleapi"
)

func Normalize(err error) error {
	if err == nil {
		return nil
	}

	if apiErr, ok := err.(*internalErrors.APIError); ok {
		return apiErr
	}

	var gErr *googleapi.Error
	if stderrors.As(err, &gErr) {
		switch gErr.Code {
		case 401:
			return internalErrors.NewGoogleReauthorizationRequiredError("Googleアカウントの再認可が必要です")
		case 403:
			return internalErrors.NewGoogleReauthorizationRequiredError("Googleアカウントの再認可が必要です")
		case 404:
			return internalErrors.NewNotFoundError("リソースが見つかりません")
		default:
			return internalErrors.NewInternalError("Google APIエラーが発生しました")
		}
	}

	return internalErrors.NewInternalError("予期せぬエラーが発生しました")
}
