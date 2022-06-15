package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	hh "github.com/0n1shi/http-honeypot"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

var version = "unknown" // overwritten by goreleaser

func main() {
	(&cli.App{
		Name:  "http honeypot",
		Usage: "http server working as honeypot",
		Commands: []*cli.Command{
			{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   "start honeyport http server",
				Action:  runServer,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "config",
						Aliases:  []string{"conf"},
						Usage:    "path of config file",
						Required: true,
					},
				},
			},
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "show version",
				Action:  showVersion,
			},
		},
	}).Run(os.Args)
}

func runServer(c *cli.Context) error {
	log.SetFlags(log.Llongfile | log.Ldate | log.Ltime)
	log.SetPrefix("[http honeypot]")

	log.Println("loading config ...")
	var conf hh.Config
	content, err := os.ReadFile(c.String("config"))
	if err != nil {
		return errors.WithStack(err)
	}
	if err := yaml.Unmarshal(content, &conf); err != nil {
		return errors.WithStack(err)
	}

	log.Println("setting up database ...")
	var repo hh.RequestRepository
	if conf.MySQL.Hostname != "" {
		repo, err = hh.NewMySQLRequestRepository(&conf.MySQL)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	rootController := hh.NewRootController(repo, &conf.Web)

	http.HandleFunc("/", rootController.HandlerAny)

	log.Printf("starting server ... :%d\n", conf.Web.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", conf.Web.Port), nil); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func showVersion(c *cli.Context) error {
	fmt.Println(version)
	return nil
}
