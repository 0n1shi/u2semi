package u2semi

type Request struct {
	Method  string
	URL     string
	Proto   string
	Headers map[string]string
	Body    string
	IPFrom  string
	IPTo    string
}

type RequestRepository interface {
	Create(req *Request) error
}
