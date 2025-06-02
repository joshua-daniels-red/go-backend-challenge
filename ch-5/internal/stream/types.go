package stream

type Event struct {
	Domain string `json:"domain"`
	Title  string `json:"title"`
	User   string `json:"user"`
}
