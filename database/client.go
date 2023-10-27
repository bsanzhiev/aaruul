package database

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"

	"github.com/bsanzhiev/tsurhai/models"
)

var Instance *gorm.DB
var dbError error

func Connect() *gorm.DB {
	Instance, dbError = gorm.Open(postgres.Open(fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=disable TimeZone=%s",
		viper.GetString("DB_HOST"),
		viper.GetInt("DB_PORT"),
		viper.GetString("DB_USER"),
		viper.GetString("DB_NAME"),
		viper.GetString("DB_PASSWORD"),
		viper.GetString("DB_TIMEZONE"),
	)), &gorm.Config{
		PrepareStmt:                              true,
		SkipDefaultTransaction:                   true,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if dbError != nil {
		panic(dbError.Error())
	}
	log.Println("Connected to database!")

	return Instance
}

func Migrate() error {
	err := Instance.AutoMigrate(
		&models.User{},
	)
	if err != nil {
		log.Println("Database Migration Failed:", err)
		return err
	}
	log.Println("Database Migration Completed!")
	return nil
}
