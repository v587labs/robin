package gorm

import "gorm.io/gorm"

const (
	SqlLite  = "sqlite"
	Postgres = "postgres"

	GroupDefault = "default"
)

type Client struct {
	db  *gorm.DB
	cfg Config
}
