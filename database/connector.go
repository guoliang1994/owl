package database

import "gorm.io/gorm"

type Connector interface {
	Open(dsn string, cfg *gorm.Config) (*gorm.DB, error)
	Options() *Options
}
