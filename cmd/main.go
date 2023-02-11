package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/0n1shi/u2semi"
	"github.com/0n1shi/u2semi/repository/mysql"
	"github.com/0n1shi/u2semi/repository/none"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

var version = "unknown" // overwritten by goreleaser

func main() {
	(&cli.App{
		Name:  "U2semi",
		Usage: "A honeypot working as a HTTP server ",
		Commands: []*cli.Command{
			{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   "Start HTTP server",
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
				Usage:   "Show version",
				Action:  showVersion,
			},
		},
	}).Run(os.Args)
}

func runServer(c *cli.Context) error {
	log.SetFlags(log.Llongfile | log.Ldate | log.Ltime)
	log.SetPrefix("[http honeypot]")

	log.Println("loading config ...")
	var conf u2semi.Conf
	content, err := os.ReadFile(c.String("config"))
	if err != nil {
		return errors.WithStack(err)
	}
	if err := yaml.Unmarshal(content, &conf); err != nil {
		return errors.WithStack(err)
	}

	log.Println("setting up repository ...")

	var repo u2semi.RequestRepository
	switch conf.Repo.Type {
	case u2semi.RepoTypeNone:
		repo = none.NewNoneRepository()
	case u2semi.RepoTypeMySQL:
		repo, err = mysql.NewMySQLRequestRepository(&conf.Repo.MySQL)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	rootController := u2semi.NewRootController(repo, &conf.Web)
	http.HandleFunc("/hello", rootController.HandlerAny)

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
