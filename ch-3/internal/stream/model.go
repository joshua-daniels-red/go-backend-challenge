package stream

type ChangeEvent struct {
	Type      string `json:"type"`
	Title     string `json:"title"`
	User      string `json:"user"`
	Bot       bool   `json:"bot"`
	ServerURL string `json:"server_url"`
}