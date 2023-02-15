package u2semi

import (
	"fmt"
	"io/ioutil"
	"log"
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
	log.Println("received a http request")
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
	if err := c.repo.Create(&req); err != nil {
		log.Println(err)
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
			if !strings.HasSuffix(r.RequestURI, "/") {
				w.Header().Set("Location", fmt.Sprintf("%s%s", r.URL, "/"))
				w.WriteHeader(http.StatusMovedPermanently)
				return
			}

			// list files
			dirListPage := DirListPageTemplate{}
			dirListPage.Dir = r.RequestURI[:len(r.RequestURI)-1]
			dirListPage.ParentDir = filepath.Dir(dirListPage.Dir)
			files, err := ioutil.ReadDir(localContentDirPath)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			for _, file := range files {
				dirListPage.Files = append(dirListPage.Files, file.Name())
			}
			t, err := template.ParseFiles(c.conf.DirListTemplate)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if err := t.Execute(w, dirListPage); err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		}

		// return file content
		content, err := ioutil.ReadFile(localContentDirPath)
		if err != nil {
			log.Fatalln(err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(content)
		return
	}

	// content from config file
	if content, ok := c.conf.Contents[r.URL.Path]; ok {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(content.Body))
		return
	}

	// content not found
	w.WriteHeader(http.StatusOK)
}
