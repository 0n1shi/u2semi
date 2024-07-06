package postgres

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PostgresRequest struct {
	gorm.Model
	Method  string `gorm:"size:256"`
	URL     string
	Proto   string `gorm:"size:256"`
	Headers datatypes.JSON
	Body    string
	IPFrom  string `gorm:"size:256"`
	IPTo    string `gorm:"size:256"`
}

func (model *PostgresRequest) TableName() string {
	return "http_requests"
}

