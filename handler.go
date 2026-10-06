package u2semi

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

// maxBodySize は保存するリクエストボディの最大バイト数（超過分は切り捨てる）
const maxBodySize = 10 << 20 // 10 MiB

type RootController struct {
	repo RequestRepository
	conf *WebConf
}

func NewRootController(repo RequestRepository, conf *WebConf) *RootController {
	return &RootController{repo: repo, conf: conf}
}

type DirListPageTemplate struct {
	ParentDir string
	Dir       string
	Files     []string
}

func (c *RootController) HandlerAny(w http.ResponseWriter, r *http.Request) {
	slog.Info("received a http request")
	req := Request{ReceivedAt: time.Now()}

	// start line
	fmt.Printf("%s %s %s\n", r.Method, r.RequestURI, r.Proto)
	req.Method = r.Method
	req.URL = r.RequestURI
	req.Proto = r.Proto
	req.IPFrom = hostOf(r.RemoteAddr)
	req.IPTo = ""
	if localAddr, ok := r.Context().Value(http.LocalAddrContextKey).(net.Addr); ok {
		req.IPTo = hostOf(localAddr.String())
	}

	// http headers
	req.Headers = make(map[string]string)
	// Host ヘッダは net/http が r.Header から取り除いて r.Host に移すため、明示的に記録する
	if r.Host != "" {
		fmt.Printf("Host: %s\n", r.Host)
		req.Headers["Host"] = r.Host
	}
	for k, v := range r.Header {
		val := strings.Join(v, " ")
		fmt.Printf("%s: %s\n", k, val)
		req.Headers[k] = val
	}

	// request body
	body, _ := io.ReadAll(io.LimitReader(r.Body, maxBodySize))
	fmt.Printf("\n%s\n", string(body))
	req.Body = string(body)

	// save request
	if err := c.repo.Save(&req); err != nil {
		slog.Error("failed to save request", "message", err.Error())
	}

	// make response header
	for _, h := range c.conf.Headers {
		w.Header().Set(h.Key, h.Value)
	}

	// content from file system
	localContentDirPath, ok := c.localPath(r.URL.Path)
	if stat, err := os.Stat(localContentDirPath); ok && err == nil { // file or directory exists
		// directory listing
		if stat.IsDir() {
			// redirect to a uri which ends with "/"
			if !strings.HasSuffix(r.URL.Path, "/") {
				u := *r.URL
				u.Path += "/"
				u.RawPath = ""
				w.Header().Set("Location", u.String())
				w.WriteHeader(http.StatusMovedPermanently)
				return
			}

			// list files
			dirListPage := DirListPageTemplate{}
			dirListPage.Dir = r.URL.Path
			if len(r.URL.Path) > 1 {
				dirListPage.Dir = r.URL.Path[:len(r.URL.Path)-1]
			}
			dirListPage.ParentDir = filepath.Dir(dirListPage.Dir)
			files, err := os.ReadDir(localContentDirPath)
			if err != nil {
				slog.Error("failed to read directory", "message", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			for _, file := range files {
				dirListPage.Files = append(dirListPage.Files, file.Name())
			}
			t, err := template.ParseFiles(c.conf.DirListTemplate)
			if err != nil {
				slog.Error("failed to parse template", "message", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// 途中まで書き出した後にエラーになるのを避けるため、一旦バッファに描画する
			var buf bytes.Buffer
			if err := t.Execute(&buf, dirListPage); err != nil {
				slog.Error("failed to execute template", "message", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			if _, err := buf.WriteTo(w); err != nil {
				slog.Error("failed to write response", "message", err.Error())
			}
			return
		}

		// return file content
		content, err := os.ReadFile(localContentDirPath)
		if err != nil {
			slog.Error("failed to read file", "message", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		if _, err = w.Write(content); err != nil {
			slog.Error("failed to write response", "message", err.Error())
		}
		return
	}

	// content from config file
	if content, ok := c.findContent(r.URL); ok {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(content.Body)); err != nil {
			slog.Error("failed to write response", "message", err.Error())
		}
		return
	}

	// content not found
	w.WriteHeader(http.StatusOK)
}

// localPath は URL のパスをコンテンツディレクトリ配下のローカルパスに変換する。
// コンテンツディレクトリが未設定の場合は false を返す。
// パスは正規化されるため、".." でコンテンツディレクトリの外に出ることはできない。
func (c *RootController) localPath(urlPath string) (string, bool) {
	if c.conf.ContentDir == "" {
		return "", false
	}
	cleaned := path.Clean("/" + urlPath)
	return filepath.Join(c.conf.ContentDir, filepath.FromSlash(cleaned)), true
}

// findContent は設定ファイルの contents からレスポンスを探す。
// クエリ付きのキー（例: "/search?q=x"）を優先し、なければパスのみのキーで探す。
func (c *RootController) findContent(u *url.URL) (*WebContent, bool) {
	if u.RawQuery != "" {
		if content, ok := c.conf.Contents[u.Path+"?"+u.RawQuery]; ok {
			return content, true
		}
	}
	content, ok := c.conf.Contents[u.Path]
	return content, ok
}

// hostOf は "host:port" 形式のアドレスからホスト部分を取り出す（IPv6 にも対応）
func hostOf(addr string) string {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return addr
	}
	return host
}
