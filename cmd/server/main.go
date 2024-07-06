package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/0n1shi/u2semi"
	"github.com/0n1shi/u2semi/repository"
	"github.com/urfave/cli/v2"
)

var version = "unknown" // overwritten by goreleaser

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})))
	app := &cli.App{
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
	}
	if err := app.Run(os.Args); err != nil {
		slog.Error("failed to run app", "message", err.Error())
		os.Exit(1)
	}
}

func runServer(c *cli.Context) error {
	slog.Info("Loading config file", "path", c.String("config"))
	conf, err := u2semi.LoadConf(c.String("config"))
	if err != nil {
		return err
	}

	repo, err := repository.NewRequestRepository(conf.Repo.DSN)
	if err != nil {
		return err
	}

	rootController := u2semi.NewRootController(repo, &conf.Web)
	http.HandleFunc("/", rootController.HandlerAny)

	slog.Info("starting server ...", "port", conf.Web.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", conf.Web.Port), nil); err != nil {
		return err
	}
	return nil
}

func showVersion(c *cli.Context) error {
	fmt.Println(version)
	return nil
}
