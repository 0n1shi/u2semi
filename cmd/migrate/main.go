package main

import (
	"log/slog"
	"os"

	"github.com/0n1shi/u2semi"
	"github.com/0n1shi/u2semi/repository"
	"gopkg.in/yaml.v2"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})))

	if len(os.Args) < 2 {
		slog.Error("no config file specified")
		os.Exit(1)
	}

	var conf u2semi.Conf
	confContent, err := os.ReadFile(os.Args[1])
	if err != nil {
		slog.Error("failed to read config file", "message", err.Error())
		os.Exit(1)
	}
	if err := yaml.Unmarshal(confContent, &conf); err != nil {
		slog.Error("failed to unmarshal config file", "message", err.Error())
		os.Exit(1)
	}

	reqRepo, err := repository.NewRequestRepository(conf.Repo.DSN)
	if err != nil {
		slog.Error("failed to create request repository", "message", err.Error())
		os.Exit(1)
	}

	if err := reqRepo.Migrate(); err != nil {
		slog.Error("failed to migrate request repository", "message", err.Error())
		os.Exit(1)
	}
}
