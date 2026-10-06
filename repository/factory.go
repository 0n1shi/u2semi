package repository

import (
	"github.com/0n1shi/u2semi"
	"github.com/0n1shi/u2semi/repository/file"
	"github.com/0n1shi/u2semi/repository/multi"
	"github.com/0n1shi/u2semi/repository/none"
	"github.com/0n1shi/u2semi/repository/postgres"
)

// NewRequestRepository は設定に応じた保存先を返す。
// dsn と file の両方が設定されている場合は両方に保存する。
func NewRequestRepository(conf *u2semi.RepoConf) (u2semi.RequestRepository, error) {
	var repos []u2semi.RequestRepository
	if conf.DSN != "" {
		repo, err := postgres.NewPostgresRequestRepository(conf.DSN)
		if err != nil {
			return nil, err
		}
		repos = append(repos, repo)
	}
	if conf.File != "" {
		repo, err := file.NewFileRequestRepository(conf.File)
		if err != nil {
			return nil, err
		}
		repos = append(repos, repo)
	}

	switch len(repos) {
	case 0:
		return none.NewNoneRepository(), nil
	case 1:
		return repos[0], nil
	default:
		return multi.NewMultiRequestRepository(repos...), nil
	}
}
