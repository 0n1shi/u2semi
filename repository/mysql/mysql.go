package mysql

import (
	"encoding/json"
	"fmt"

	httphoneypot "github.com/0n1shi/http-honeypot"
	"github.com/pkg/errors"
	"gorm.io/datatypes"
	_mysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type RequestMySQLModel struct {
	gorm.Model
	Method  string `gorm:"size:256"`
	URL     string `gorm:"size:2048"`
	Proto   string `gorm:"size:256"`
	Headers datatypes.JSON
	Body    string
	IPFrom  string `gorm:"size:256"`
	IPTo    string `gorm:"size:256"`
}

func (model *RequestMySQLModel) TableName() string {
	return "http_requests"
}

var _ httphoneypot.RequestRepository = (*MySQLRequestRepository)(nil)

type MySQLRequestRepository struct {
	db *gorm.DB
}

func NewMySQLRequestRepository(conf *httphoneypot.MySQLConfig) (*MySQLRequestRepository, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.Username, conf.Password, conf.Hostname, conf.DB)
	db, err := gorm.Open(_mysql.Open(dsn), &gorm.Config{})
	return &MySQLRequestRepository{db: db}, err
}

func (repo *MySQLRequestRepository) Create(req *httphoneypot.Request) error {
	headers, err := json.Marshal(req.Headers)
	if err != nil {
		return errors.WithStack(err)
	}
	model := RequestMySQLModel{
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
