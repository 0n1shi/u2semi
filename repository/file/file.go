package file

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/0n1shi/u2semi"
)

var _ u2semi.RequestRepository = (*FileRequestRepository)(nil)

// record はファイルに 1 行ずつ書き出す JSON の形式
type record struct {
	ReceivedAt time.Time         `json:"received_at"`
	Method     string            `json:"method"`
	URL        string            `json:"url"`
	Proto      string            `json:"proto"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
	// BodyBase64 はボディが UTF-8 として不正な場合のみ、生のバイト列を保持する
	BodyBase64 string `json:"body_base64,omitempty"`
	IPFrom     string `json:"ip_from"`
	IPTo       string `json:"ip_to"`
}

// FileRequestRepository はリクエストを JSON Lines 形式でファイルに追記する
type FileRequestRepository struct {
	mu   sync.Mutex
	file *os.File
}

// NewFileRequestRepository は path のファイルを追記モードで開く（なければ親ディレクトリごと作成する）
func NewFileRequestRepository(path string) (*FileRequestRepository, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("failed to create directory for %s: %w", path, err)
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", path, err)
	}
	return &FileRequestRepository{file: f}, nil
}

func (repo *FileRequestRepository) Save(req *u2semi.Request) error {
	rec := record{
		ReceivedAt: req.ReceivedAt,
		Method:     req.Method,
		URL:        req.URL,
		Proto:      req.Proto,
		Headers:    req.Headers,
		Body:       req.Body,
		IPFrom:     req.IPFrom,
		IPTo:       req.IPTo,
	}
	if !utf8.ValidString(req.Body) {
		rec.BodyBase64 = base64.StdEncoding.EncodeToString([]byte(req.Body))
	}
	line, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	line = append(line, '\n')

	// 並行リクエストで行が混ざらないよう、1 行を 1 回の Write で書き込む
	repo.mu.Lock()
	defer repo.mu.Unlock()
	_, err = repo.file.Write(line)
	return err
}

func (repo *FileRequestRepository) Migrate() error {
	slog.Info("Request repository migrated (file repo, nothing to do)")
	return nil
}

// Close はファイルを閉じる
func (repo *FileRequestRepository) Close() error {
	return repo.file.Close()
}
