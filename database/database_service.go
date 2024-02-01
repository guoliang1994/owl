package database

import (
	"fmt"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"owl"
	"owl/contract"
	"strconv"
	"sync"
	"time"
)

var (
	connections = make(map[string]*gorm.DB, 10)
	lock        sync.Mutex
)

// 能力
// 支持连接多个不同的数据库
// 内部维护数据库连接池
// 使用依赖注入时，返回新的连接

type Options struct {
	contract.ServerConfig
	Driver       string `json:"driver"` // 数据库类型
	Database     string `json:"database"`
	Schema       string `json:"schema"`
	Charset      string `json:"charset"`
	Query        string `json:"query"`
	MaxIdleConns int    `json:"max-idle-conns"`
	MaxConns     int    `json:"max-conns"`
}

func NewOption(stage *owl.Stage) (opt *Options) {
	err := jsoniter.UnmarshalFromString(stage.ConfManager.GetConfig("db").ToString(), &opt)
	if err != nil {
		return nil
	}
	return opt
}

type DatabaseService struct {
	stage    *owl.Stage
	l        contract.Logger
	dsn      string
	opt      *Options
	db       *gorm.DB
	dbGetter Connector
}

func NewDatabaseService(stage *owl.Stage, dbGetter Connector) *DatabaseService {
	opt := dbGetter.Options()

	i := &DatabaseService{
		stage:    stage,
		opt:      opt,
		dbGetter: dbGetter,
	}

	i.Boot()
	return i
}

// Boot 打开连接
func (i *DatabaseService) Boot() {
	lock.Lock()
	defer lock.Unlock()

	cacheDb, ok := connections[i.dsn] // 申请连接过

	if !ok || cacheDb == nil {
		i.dsn = i.getDsnFromCfg(i.opt)

		i.watchCfgFileChange()

		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second,   // 慢 SQL 阈值
				LogLevel:      logger.Silent, // Log level
				Colorful:      false,         // 禁用彩色打印
			},
		)

		gormCfg := &gorm.Config{
			PrepareStmt:                              false,
			DisableForeignKeyConstraintWhenMigrating: true,
			NamingStrategy: schema.NamingStrategy{
				//单数表名
				SingularTable: true,
			},
			Logger:                 newLogger,
			SkipDefaultTransaction: true,
		}

		var openDb *gorm.DB
		var err error
		openDb, err = i.dbGetter.Open(i.dsn, gormCfg)

		if err != nil {

			panic("数据库连接失败，请检查数据库是否启动，配置是否错误" + err.Error())
		}
		sqlDB, err := openDb.DB()

		// SetMaxIdleConns 设置空闲连接池中连接的最大数量
		sqlDB.SetMaxIdleConns(i.opt.MaxIdleConns)

		// SetMaxOpenConns 设置打开数据库连接的最大数量。
		// <= 数据库配置的连接数量
		sqlDB.SetMaxOpenConns(i.opt.MaxConns)

		// SetConnMaxLifetime 设置了连接可复用的最大时间。
		sqlDB.SetConnMaxLifetime(time.Hour)
		if err != nil {
			return
		}

		i.db = openDb
		connections[i.dsn] = openDb
	}

	return
}

func (i *DatabaseService) Get() *gorm.DB {
	return i.db
}

func (i *DatabaseService) getDsnFromCfg(opt *Options) string {
	var dsn string

	if opt.Driver == "mysql" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&%s",
			opt.Username,
			opt.Password,
			opt.Host,
			opt.Port,
			opt.Database,
			opt.Charset,
			opt.Query,
		)

	} else if opt.Driver == "sqlite" {
		dsn = opt.Host
	}
	return dsn
}

func (i *DatabaseService) watchCfgFileChange() {
	go func() {
		for {
			select {
			case file := <-owl.CfgChangeNotify[i.opt.AbsPath]:
				lock.Lock()
				i.l.Info("数据库重连" + file)
				delete(connections, i.dsn)
				lock.Unlock()
				return
			}
		}
	}()
}

func Paginate(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page, _ := strconv.Atoi(c.Query("page"))
		if page == 0 {
			page = 1
		}

		pageSize, _ := strconv.Atoi(c.Query("pageSize"))
		switch {
		case pageSize > 100:
			pageSize = 1000
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
