package u2semi

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
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

func TestHandlerAny_BodyIsTruncated(t *testing.T) {
	srv, repo := newTestServer(t, &WebConf{})

	body := strings.Repeat("a", maxBodySize+100)
	resp, err := http.Post(srv.URL+"/", "text/plain", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if got := len(repo.last(t).Body); got != maxBodySize {
		t.Errorf("len(Body) = %d, want %d", got, maxBodySize)
	}
}

// get は生のリクエストパスで GET し、ステータス・Location・ボディを返す
func get(t *testing.T, srv *httptest.Server, rawPath string) (int, string, string) {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.URL.Opaque = rawPath // エンコードを変えずにそのまま送る
	client := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, resp.Header.Get("Location"), string(b)
}

func newContentDir(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	dir := filepath.Join(root, "content")
	if err := os.MkdirAll(filepath.Join(dir, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "sub", "a.txt"), []byte("A"), 0o644); err != nil {
		t.Fatal(err)
	}
	// コンテンツディレクトリの外にあるファイル（配信されてはいけない）
	if err := os.WriteFile(filepath.Join(root, "secret.txt"), []byte("SECRET"), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestHandlerAny_ContentFile(t *testing.T) {
	srv, _ := newTestServer(t, &WebConf{
		ContentDir:      newContentDir(t),
		DirListTemplate: "template/directory_listing.html",
	})

	tests := []struct {
		path     string
		status   int
		location string
		body     string
	}{
		{path: "/sub/a.txt", status: 200, body: "A"},
		{path: "/sub/a.txt?x=1", status: 200, body: "A"},
		{path: "/sub?x=1", status: 301, location: "/sub/?x=1"},
		{path: "/%2e%2e/secret.txt", status: 200, body: ""},
		{path: "/sub/..%2f..%2fsecret.txt", status: 200, body: ""},
	}
	for _, tt := range tests {
		status, location, body := get(t, srv, tt.path)
		if status != tt.status || location != tt.location || body != tt.body {
			t.Errorf("GET %s = (%d, %q, %q), want (%d, %q, %q)",
				tt.path, status, location, body, tt.status, tt.location, tt.body)
		}
	}
}

func TestHandlerAny_NoContentDirDoesNotServeFilesystem(t *testing.T) {
	srv, _ := newTestServer(t, &WebConf{})

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	_, _, body := get(t, srv, filepath.ToSlash(filepath.Join(wd, "go.mod")))
	if body != "" {
		t.Errorf("served a local file without content_directory: %q", body)
	}
}
