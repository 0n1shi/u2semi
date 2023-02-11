package none

import (
	"github.com/0n1shi/u2semi"
	"gorm.io/gorm"
)

var _ u2semi.RequestRepository = (*NoneRepository)(nil)

type NoneRepository struct {
	db *gorm.DB
}

func NewNoneRepository() *NoneRepository {
	return &NoneRepository{}
}

func (repo *NoneRepository) Create(req *u2semi.Request) error {
	return nil
}
