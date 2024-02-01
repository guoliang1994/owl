package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"owl/contract"
)

type sqliteGetter struct {
	opt *Options
}

func NewSqliteGetter(opt *Options) *sqliteGetter {
	return &sqliteGetter{opt}
}
func (i *sqliteGetter) Open(dsn string, cfg *gorm.Config) (*gorm.DB, error) {
	openDb, err := gorm.Open(sqlite.Open(dsn), cfg)
	return openDb, err

}
func (i *sqliteGetter) Options() *Options {
	if i.opt != nil {
		return i.opt
	}

	return &Options{
		ServerConfig: contract.ServerConfig{
			Host:     "owlSqlite.db",
			Username: "root",
			Password: "root",
		},
		Database:     "main",
		Driver:       "sqlite",
		MaxIdleConns: 10,
		MaxConns:     100,
	}
}
