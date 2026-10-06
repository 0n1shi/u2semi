package multi

import (
	"errors"

	"github.com/0n1shi/u2semi"
)

var _ u2semi.RequestRepository = (*MultiRequestRepository)(nil)

// MultiRequestRepository は複数のリポジトリに同じリクエストを保存する
type MultiRequestRepository struct {
	repos []u2semi.RequestRepository
}

func NewMultiRequestRepository(repos ...u2semi.RequestRepository) *MultiRequestRepository {
	return &MultiRequestRepository{repos: repos}
}

// Save は全リポジトリに保存を試み、失敗したものがあればまとめてエラーを返す
func (repo *MultiRequestRepository) Save(req *u2semi.Request) error {
	var errs []error
	for _, r := range repo.repos {
		if err := r.Save(req); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (repo *MultiRequestRepository) Migrate() error {
	var errs []error
	for _, r := range repo.repos {
		if err := r.Migrate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
