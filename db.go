package httphoneypot

type RequestModel struct {
	Method string            `json:"method"`
	URL    string            `json:"url"`
	Proto  string            `json:"version"`
	Header map[string]string `json:"headers"`
	Body   string            `json:"body"`
}
