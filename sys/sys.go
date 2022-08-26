package sys

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

var (
	Conn *gorm.DB
	L    string
	U    string
	P    string
)

func ConnectMySQL() error {
	logout, err := os.OpenFile("./db.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=%s&parseTime=True&loc=Local", U, P, L, U, "utf8mb4")
	Conn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(logout, "\r\n", log.Lshortfile|log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Error,
				Colorful:      false,
			},
		),
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",
			SingularTable: true,
		},
		PrepareStmt: true, // 预编译模式
	})
	if err != nil {
		return err
	}
	if conn, ok := Conn.ConnPool.(*sql.DB); ok {
		conn.SetMaxIdleConns(200)
		conn.SetMaxOpenConns(199)
		conn.SetConnMaxLifetime(time.Hour)
		conn.SetConnMaxIdleTime(time.Hour)
	}
	return nil
}
