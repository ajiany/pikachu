package config

import (
	"fmt"
	"os"
	"time"

	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	gormtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorm.io/gorm.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB mysql服务初始化
func InitDB() {
	var err error
	sqltrace.Register("mysql", &mysqldriver.MySQLDriver{}, sqltrace.WithServiceName(os.Getenv("DD_SERVICE")))
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", Cfg.DBUsername, Cfg.DBPassword, Cfg.DBHostname, Cfg.DBPort, Cfg.DBDatabase)
	sqlDb, err := sqltrace.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	DB, err = gormtrace.Open(mysql.New(mysql.Config{Conn: sqlDb}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	pool, err := DB.DB()
	if err != nil {
		panic(err)
	}
	pool.SetMaxIdleConns(10)
	pool.SetMaxOpenConns(100)
	pool.SetConnMaxLifetime(time.Hour)

	logrus.Info("init db success.")
}
