package httphoneypot

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type RootController struct {
	repo RequestRepository
	conf *WebConfig
}

func NewRootController(repo RequestRepository, conf *WebConfig) *RootController {
	return &RootController{repo: repo, conf: conf}
}

func (c *RootController) HandlerAny(w http.ResponseWriter, r *http.Request) {
	log.Println("received a http request")
	req := Request{}

	// start line
	fmt.Printf("%s %s %s\n", r.Method, r.RequestURI, r.Proto)
	req.Method = r.Method
	req.URL = r.RequestURI
	req.Proto = r.Proto
	req.IP = r.RemoteAddr

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
	if c.repo != nil {
		if err := c.repo.Create(&req); err != nil {
			log.Println(err)
		}
	}

	// make response
	for _, h := range c.conf.Headers {
		w.Header().Set(h.Key, h.Value)
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message": "hello world"}`)) // TODO
}
