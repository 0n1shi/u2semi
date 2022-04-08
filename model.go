package httphoneypot

type Request struct {
	Method string
	URL    string
	Proto  string // "HTTP/1.0"
	Header map[string][]string
	Body   string
}

type StartLine struct {
	Method  string
	URL     string
	Version string
}
