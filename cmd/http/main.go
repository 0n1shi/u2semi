package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	hh "github.com/0n1shi/http-honeypot"
	"github.com/urfave/cli/v2"
)

func main() {
	(&cli.App{
		Name:  "http honeypot",
		Usage: "http server working as honeypot",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "conf",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "serve",
				Aliases: []string{"s"},
				Usage:   "start honeyport http server",
				Action:  serve,
			},
		},
	}).Run(os.Args)

}

func serve(c *cli.Context) error {
	log.SetFlags(log.Llongfile | log.Ldate | log.Ltime)
	log.SetPrefix("[http honeypot]")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("received a http request")
		req := hh.Request{}
		fmt.Printf("%s %s %s\n", r.Method, r.RequestURI, r.Proto)
		req.Method = r.Method
		req.URL = r.RequestURI
		req.Proto = r.Proto
		for k, v := range r.Header {
			val := strings.Join(v, " ")
			fmt.Printf("%s: %s\n", k, val)
			req.Header[k] = val
		}
		body, _ := ioutil.ReadAll(r.Body)
		fmt.Printf("\n%s\n", string(body))
		req.Body = string(body)
	})

	http.ListenAndServe(":8080", nil)
	return nil
}
