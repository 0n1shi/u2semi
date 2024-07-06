package none

import (
	"github.com/0n1shi/u2semi"
	"log/slog"
)

var _ u2semi.RequestRepository = (*NoneRepository)(nil)

func NewNoneRepository() *NoneRepository {
	return &NoneRepository{}
}

type NoneRepository struct {
}

func (repo *NoneRepository) Save(req *u2semi.Request) error {
	slog.Info("Request saved (none repo)", "request", req)
	return nil
}

func (repo *NoneRepository) Migrate() error {
	slog.Info("Request repository migrated (none repo)")
	return nil
}
