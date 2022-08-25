package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDB() (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?timeout=5s&readTimeout=5s",
		viper.GetString("data.database.user"),
		viper.GetString("data.database.password"),
		viper.GetString("data.database.host"),
		viper.GetString("data.database.port"),
		viper.GetString("data.database.database"))
	return sql.Open(viper.GetString("data.database.driver"), dsn)
}
func NewGormDB(sqlDB *sql.DB) *gorm.DB {
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		log.Printf("opening connection to mysql: %v is failed", err)
		return nil
	}

	//if err = db.AutoMigrate(&v1alpha1.Demo{}); err != nil {
	//	log.Printf("database auto migrate is failed: %v", err)
	//	return nil
	//}
	return db
}
