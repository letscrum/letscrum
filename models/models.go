package models

import (
	"fmt"
	"gorm.io/gorm/schema"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/letscrum/letscrum/pkg/settings"
)

var db *gorm.DB

type Model struct {
	ID int64 `json:"id" gorm:"primary_key"`
}

// Setup initializes the database instance
func Setup() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		settings.DatabaseSetting.User,
		settings.DatabaseSetting.Password,
		settings.DatabaseSetting.Host,
		settings.DatabaseSetting.Name)
	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			TablePrefix:   settings.DatabaseSetting.TablePrefix,
		},
	})
	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}
	sqlDB, errDB := conn.DB()
	if errDB != nil {
		log.Fatalf("models.Setup err: %v", errDB)
	}
	sqlDB.SetMaxOpenConns(1000)
	sqlDB.SetConnMaxLifetime(10)
	sqlDB.SetMaxIdleConns(10)

	db = conn
}
