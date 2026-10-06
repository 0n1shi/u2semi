package u2semi

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

// fakeRepo は保存されたリクエストを記録するテスト用リポジトリ
type fakeRepo struct {
	mu   sync.Mutex
	reqs []*Request
}

func (f *fakeRepo) Save(req *Request) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.reqs = append(f.reqs, req)
	return nil
}

func (f *fakeRepo) Migrate() error { return nil }

func (f *fakeRepo) last(t *testing.T) *Request {
	t.Helper()
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.reqs) == 0 {
		t.Fatal("no request saved")
	}
	return f.reqs[len(f.reqs)-1]
}

func newTestServer(t *testing.T, conf *WebConf) (*httptest.Server, *fakeRepo) {
	t.Helper()
	repo := &fakeRepo{}
	srv := httptest.NewServer(http.HandlerFunc(NewRootController(repo, conf).HandlerAny))
	t.Cleanup(srv.Close)
	return srv, repo
}

func TestHandlerAny_IPAddresses(t *testing.T) {
	srv, repo := newTestServer(t, &WebConf{})

	resp, err := http.Get(srv.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	req := repo.last(t)
	if req.IPFrom != "127.0.0.1" {
		t.Errorf("IPFrom = %q, want 127.0.0.1", req.IPFrom)
	}
	if req.IPTo != "127.0.0.1" {
		t.Errorf("IPTo = %q, want 127.0.0.1", req.IPTo)
	}
}

func TestHostOf(t *testing.T) {
	tests := map[string]string{
		"192.0.2.1:1234":     "192.0.2.1",
		"[2001:db8::1]:8080": "2001:db8::1",
		"invalid":            "invalid",
	}
	for in, want := range tests {
		if got := hostOf(in); got != want {
			t.Errorf("hostOf(%q) = %q, want %q", in, got, want)
		}
	}
}
