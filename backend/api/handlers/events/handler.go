package events

type Handler struct {
	googleCalendarUsecase GoogleCalendarUsecase
	draftUsecase          DraftUsecase
	detailUsecase         DetailUsecase
	agendaUsecase         AgendaUsecase
	confirmationUsecase   ConfirmationUsecase
}

func NewHandler(
	googleCalendarUsecase GoogleCalendarUsecase,
	draftUsecase DraftUsecase,
	detailUsecase DetailUsecase,
	agendaUsecase AgendaUsecase,
	confirmationUsecase ConfirmationUsecase,
) *Handler {
	return &Handler{
		googleCalendarUsecase: googleCalendarUsecase,
		draftUsecase:          draftUsecase,
		detailUsecase:         detailUsecase,
		agendaUsecase:         agendaUsecase,
		confirmationUsecase:   confirmationUsecase,
	}
}

var extractErrorMessage = "ユーザー情報確認時にエラーが発生しました。"
