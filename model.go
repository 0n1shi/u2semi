package httphoneypot

type Request struct {
	StartLine *StartLine
}

type StartLine struct {
	Method  string
	URL     string
	Version string
}
