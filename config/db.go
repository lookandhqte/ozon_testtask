package config

import (
	"ozon_testTask/graph/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDb() (*gorm.DB, error) {
	var err error
	dsn := "host=localhost user=postgres password=postgres dbname=ozon_test port=5432"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&model.Comment{}, &model.Post{}, &model.Query{}, &model.Subscription{})

	return db, nil
}
