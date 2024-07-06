package u2semi

import (
	"fmt"
	"io/ioutil"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

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
	req := Request{}

	// start line
	fmt.Printf("%s %s %s\n", r.Method, r.RequestURI, r.Proto)
	req.Method = r.Method
	req.URL = r.RequestURI
	req.Proto = r.Proto
	req.IPFrom = strings.Split(r.RemoteAddr, ":")[0]
	req.IPTo = ""
	localAddr, ok := r.Context().Value(http.LocalAddrContextKey).(string)
	if ok {
		req.IPTo = localAddr
	}

	// http headers
	req.Headers = make(map[string]string)
	for k, v := range r.Header {
		val := strings.Join(v, " ")
		fmt.Printf("%s: %s\n", k, val)
		req.Headers[k] = val
	}

	// request body
	body, _ := ioutil.ReadAll(r.Body)
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
	localContentDirPath := fmt.Sprintf("%s%s", c.conf.ContentDir, r.RequestURI)
	if stat, err := os.Stat(localContentDirPath); !os.IsNotExist(err) { // directory exists
		// directory listing
		if stat.IsDir() {
			// redirect to a uri which ends with "/"
			if !strings.HasSuffix(r.URL.Path, "/") {
				w.Header().Set("Location", fmt.Sprintf("%s%s", r.URL, "/"))
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
			files, err := ioutil.ReadDir(localContentDirPath)
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
			if err := t.Execute(w, dirListPage); err != nil {
				slog.Error("failed to execute template", "message", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		}

		// return file content
		content, err := ioutil.ReadFile(localContentDirPath)
		if err != nil {
			slog.Error("failed to read file", "message", err.Error())
			os.Exit(1)
		}
		w.WriteHeader(http.StatusOK)
		if _, err = w.Write(content); err != nil {
			slog.Error("failed to write response", "message", err.Error())
			os.Exit(1)
		}
		return
	}

	// content from config file
	if content, ok := c.conf.Contents[r.URL.Path+"?"+r.URL.RawQuery]; ok {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(content.Body)); err != nil {
			slog.Error("failed to write response", "message", err.Error())
			os.Exit(1)
		}
		return
	}

	// content not found
	w.WriteHeader(http.StatusOK)
}
