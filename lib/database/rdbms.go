package database

import (
	"io"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//   driver: mysql
//   source: user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
//
//   driver: postgres
//   source: host=localhost user=username password=password dbname=dbname port=9920 sslmode=disable
//
//   driver: sqlite
//   source: /path/dbname.db
//
//   driver: sqlserver
//   source: sqlserver://username:password@localhost:9930?database=dbname
//
//   driver: clickhouse
//   source: tcp://localhost:9000?database=dbname&username=username&password=password&read_timeout=10&write_timeout=20

var db *gorm.DB

// Init database init
func InitRDBMS(logWriter io.Writer, name, source string, active, idle, idleTimeout int) (*gorm.DB, error) {
	var (
		gormDB *gorm.DB
		err    error
	)

	var newLogger logger.Interface
	if logWriter == nil {
		logWriter = os.Stdout
	}
	newLogger = logger.New(
		log.New(logWriter, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Warn, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)

	gormCfg := &gorm.Config{Logger: newLogger, SkipDefaultTransaction: true}

	switch name {
	case "mysql":
		// https://github.com/go-sql-driver/mysql
		gormDB, err = gorm.Open(mysql.Open(source), gormCfg)
	case "postgres":
		// https://github.com/go-gorm/postgres
		gormDB, err = gorm.Open(postgres.Open(source), gormCfg)
	case "sqlite":
		// github.com/mattn/go-sqlite3
		gormDB, err = gorm.Open(sqlite.Open(source), gormCfg)
		// case "sqlserver":
		// 	// github.com/denisenkom/go-mssqldb
		// 	gormDB, err = gorm.Open(sqlserver.Open(source), gormCfg)
		// case "clickhouse":
		// 	gormDB, err = gorm.Open(clickhouse.Open(source), gormCfg)
	}
	if err != nil {
		return nil, err
	}

	db = gormDB

	dd, _ := db.DB()

	dd.SetMaxOpenConns(active)
	dd.SetMaxIdleConns(idle)
	dd.SetConnMaxLifetime(time.Duration(idleTimeout) * time.Second)

	return gormDB, nil
}

func Get() *gorm.DB {
	return db
}
