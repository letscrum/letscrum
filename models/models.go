package models

import (
	"fmt"
	"github.com/golang/protobuf/ptypes/timestamp"
	"gorm.io/gorm/schema"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/letscrum/letscrum/pkg/setting"
)

var db *gorm.DB

type Model struct {
	ID        int                 `json:"id" gorm:"primary_key"`
	DeletedAt timestamp.Timestamp `json:"deleted_at"`
	CreatedAt timestamp.Timestamp `json:"created_at"`
	UpdatedAt timestamp.Timestamp `json:"updated_at"`
}

// Setup initializes the database instance
func Setup() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name)
	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			TablePrefix:   setting.DatabaseSetting.TablePrefix,
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
