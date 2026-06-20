package appmodel

type GoogleEvent struct {
	ID          string `json:"id"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Location    string `json:"location"`
	ColorID     string `json:"color"`
	Start       string `json:"start"`
	End         string `json:"end"`
}
