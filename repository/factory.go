package repository

import (
	"github.com/0n1shi/u2semi"
	"github.com/0n1shi/u2semi/repository/none"
	"github.com/0n1shi/u2semi/repository/postgres"
)

func NewRequestRepository(dsn string) (u2semi.RequestRepository, error) {
	if dsn == "" {
		return none.NewNoneRepository(), nil
	}
	return postgres.NewPostgresRequestRepository(dsn)
}
