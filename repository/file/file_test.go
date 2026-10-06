package file

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/0n1shi/u2semi"
)

func readRecords(t *testing.T, path string) []record {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var recs []record
	sc := bufio.NewScanner(f)
	sc.Buffer(nil, 1<<20)
	for sc.Scan() {
		var rec record
		if err := json.Unmarshal(sc.Bytes(), &rec); err != nil {
			t.Fatalf("invalid json line %q: %v", sc.Text(), err)
		}
		recs = append(recs, rec)
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}
	return recs
}

func TestSave(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "requests.jsonl")
	repo, err := NewFileRequestRepository(path)
	if err != nil {
		t.Fatal(err)
	}
	defer repo.Close()

	now := time.Date(2026, 10, 6, 12, 0, 0, 0, time.UTC)
	req := &u2semi.Request{
		ReceivedAt: now,
		Method:     "POST",
		URL:        "/login?x=1",
		Proto:      "HTTP/1.1",
		Headers:    map[string]string{"User-Agent": "curl/8.0"},
		Body:       "user=admin",
		IPFrom:     "192.0.2.1",
		IPTo:       "198.51.100.1",
	}
	if err := repo.Save(req); err != nil {
		t.Fatal(err)
	}

	recs := readRecords(t, path)
	if len(recs) != 1 {
		t.Fatalf("got %d records, want 1", len(recs))
	}
	got := recs[0]
	if !got.ReceivedAt.Equal(now) || got.Method != "POST" || got.URL != "/login?x=1" ||
		got.Headers["User-Agent"] != "curl/8.0" || got.Body != "user=admin" ||
		got.IPFrom != "192.0.2.1" || got.IPTo != "198.51.100.1" || got.BodyBase64 != "" {
		t.Errorf("unexpected record: %+v", got)
	}
}

func TestSave_BinaryBody(t *testing.T) {
	path := filepath.Join(t.TempDir(), "requests.jsonl")
	repo, err := NewFileRequestRepository(path)
	if err != nil {
		t.Fatal(err)
	}
	defer repo.Close()

	body := "\xff\xfe\x00bin"
	if err := repo.Save(&u2semi.Request{Body: body}); err != nil {
		t.Fatal(err)
	}

	decoded, err := base64.StdEncoding.DecodeString(readRecords(t, path)[0].BodyBase64)
	if err != nil {
		t.Fatal(err)
	}
	if string(decoded) != body {
		t.Errorf("body_base64 decodes to %q, want %q", decoded, body)
	}
}

func TestSave_AppendsAcrossReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "requests.jsonl")
	for i := 0; i < 2; i++ {
		repo, err := NewFileRequestRepository(path)
		if err != nil {
			t.Fatal(err)
		}
		if err := repo.Save(&u2semi.Request{URL: fmt.Sprintf("/%d", i)}); err != nil {
			t.Fatal(err)
		}
		repo.Close()
	}
	if got := len(readRecords(t, path)); got != 2 {
		t.Errorf("got %d records, want 2", got)
	}
}

func TestSave_Concurrent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "requests.jsonl")
	repo, err := NewFileRequestRepository(path)
	if err != nil {
		t.Fatal(err)
	}
	defer repo.Close()

	const n = 200
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if err := repo.Save(&u2semi.Request{URL: fmt.Sprintf("/%d", i)}); err != nil {
				t.Error(err)
			}
		}(i)
	}
	wg.Wait()

	if got := len(readRecords(t, path)); got != n {
		t.Errorf("got %d records, want %d", got, n)
	}
}
