package u2semi

import "time"

type Request struct {
	ReceivedAt time.Time
	Method     string
	URL        string
	Proto      string
	Headers    map[string]string
	Body       string
	IPFrom     string
	IPTo       string
}
