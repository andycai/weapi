package database

import (
	"time"

	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
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
func InitRDBMS(name, source string, active, idle, idleTimeout int) (*gorm.DB, error) {
	var (
		gormDB *gorm.DB
		err    error
	)
	switch name {
	case "mysql":
		// https://github.com/go-sql-driver/mysql
		gormDB, err = gorm.Open(mysql.Open(source), &gorm.Config{})
	case "postgres":
		// https://github.com/go-gorm/postgres
		gormDB, err = gorm.Open(postgres.Open(source), &gorm.Config{})
	case "sqlite":
		// github.com/mattn/go-sqlite3
		gormDB, err = gorm.Open(sqlite.Open(source), &gorm.Config{})
	case "sqlserver":
		// github.com/denisenkom/go-mssqldb
		gormDB, err = gorm.Open(sqlserver.Open(source), &gorm.Config{})
	case "clickhouse":
		gormDB, err = gorm.Open(clickhouse.Open(source), &gorm.Config{})
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
