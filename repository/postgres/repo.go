package postgres

import (
	"encoding/json"

	"github.com/0n1shi/u2semi"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)
var _ u2semi.RequestRepository = (*PostgresRequestRepository)(nil)

func NewPostgresRequestRepository(dsn string) (*PostgresRequestRepository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return &PostgresRequestRepository{db: db}, err
}

type PostgresRequestRepository struct {
	db *gorm.DB
}

func (repo *PostgresRequestRepository) Save(req *u2semi.Request) error {
	headers, err := json.Marshal(req.Headers)
	if err != nil {
		return err
	}
	model := PostgresRequest{
		Method:  req.Method,
		URL:     req.URL,
		Proto:   req.Proto,
		Headers: headers,
		Body:    req.Body,
		IPFrom:  req.IPFrom,
		IPTo:    req.IPTo,
	}
	return repo.db.Create(&model).Error
}

func (repo *PostgresRequestRepository) Migrate() error {
	return repo.db.AutoMigrate(&PostgresRequest{})
}
