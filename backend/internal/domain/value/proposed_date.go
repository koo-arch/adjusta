package value

type ProposedDateStatus string

const (
	ProposedDateStatusActive      ProposedDateStatus = "active"
	ProposedDateStatusConfirmed   ProposedDateStatus = "confirmed"
	ProposedDateStatusNotSelected ProposedDateStatus = "not_selected"
	ProposedDateStatusCancelled   ProposedDateStatus = "cancelled"
)
