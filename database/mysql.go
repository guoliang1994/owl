package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"owl/contract"
)

type MysqlConnector struct {
	opt *Options
}

func NewMysqlGetter(opt *Options) *MysqlConnector {
	return &MysqlConnector{opt: opt}
}
func (i *MysqlConnector) Open(dsn string, cfg *gorm.Config) (*gorm.DB, error) {
	openDb, err := gorm.Open(mysql.Open(dsn), cfg)
	return openDb, err
}

func (i *MysqlConnector) Options() *Options {
	if i.opt != nil {
		return i.opt
	}

	return i.DefaultOptions()
}

func (i *MysqlConnector) DefaultOptions() *Options {

	return &Options{
		ServerConfig: contract.ServerConfig{
			Host:     "127.0.0.1",
			Port:     3306,
			Username: "root",
			Password: "root",
		},
		Driver:       "mysql",
		Database:     "mysql",
		Charset:      "utf8mb4",
		Query:        "parseTime=True&loc=Local&timeout=3000ms",
		MaxIdleConns: 10,
		MaxConns:     100,
	}
}
