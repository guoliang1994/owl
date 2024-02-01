package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"owl/contract"
)

type pgsqlConnector struct {
	opt *Options
}

func NewPgSqlGetter(opt *Options) *pgsqlConnector {
	return &pgsqlConnector{opt}
}
func (i *pgsqlConnector) Open(dsn string, cfg *gorm.Config) (*gorm.DB, error) {
	openDb, err := gorm.Open(postgres.Open(dsn), cfg)
	return openDb, err
}

func (i *pgsqlConnector) Options() *Options {
	if i.opt != nil {
		return i.opt
	}
	return &Options{
		ServerConfig: contract.ServerConfig{
			Host:     "localhost",
			Port:     5432,
			Username: "postgres",
			Password: "postgres",
		},
		Driver:       "postgres",
		Database:     "postgres",
		Schema:       "public",
		MaxIdleConns: 10,
		MaxConns:     100,
	}
}
