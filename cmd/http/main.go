package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	hh "github.com/0n1shi/http-honeypot"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

func main() {
	(&cli.App{
		Name:  "http honeypot",
		Usage: "http server working as honeypot",
		Commands: []*cli.Command{
			{
				Name:    "serve",
				Aliases: []string{"s"},
				Usage:   "start honeyport http server",
				Action:  serve,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "config",
						Aliases:  []string{"conf"},
						Usage:    "path of config file",
						Required: true,
					},
				},
			},
		},
	}).Run(os.Args)

}

func serve(c *cli.Context) error {
	log.SetFlags(log.Llongfile | log.Ldate | log.Ltime)
	log.SetPrefix("[http honeypot]")

	var conf hh.Config
	content, err := os.ReadFile(c.String("config"))
	if err != nil {
		return errors.WithStack(err)
	}
	if err := yaml.Unmarshal(content, &conf); err != nil {
		return errors.WithStack(err)
	}

	requestRepo, err := hh.NewMySQLRequestRepository(&conf.MySQL)
	if err != nil {
		return errors.WithStack(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("received a http request")
		req := hh.Request{}

		// start line
		fmt.Printf("%s %s %s\n", r.Method, r.RequestURI, r.Proto)
		req.Method = r.Method
		req.URL = r.RequestURI
		req.Proto = r.Proto

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

		if err := requestRepo.Create(&req); err != nil {
			log.Println(err)
		}

		w.Header().Set("Server", "Apache/2.4.2 (Unix) PHP/4.2.2")
		w.WriteHeader(200)
		w.Write([]byte(`{"hello": "world"}`))
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
