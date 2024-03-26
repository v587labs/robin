package gorm

import "gorm.io/gorm"

const (
	SqlLite  = "sqlite"
	Postgres = "postgres"

	GroupDefault = "default"
)

type Client struct {
	*gorm.DB
	cfg Config
}
