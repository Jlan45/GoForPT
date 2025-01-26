package database

import (
	"GoForPT/model"
	"GoForPT/pkg/cfg"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai", cfg.Cfg.Database.Host, cfg.Cfg.Database.Username, cfg.Cfg.Database.Password, cfg.Cfg.Database.Name, cfg.Cfg.Database.Port)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}
func InitDBData() error {
	//一些默认数据
	group := model.UserGroup{
		Name:  "Normal",
		Level: 1,
	}
	DB.AutoMigrate(&model.UserGroup{})
	DB.AutoMigrate(&model.User{})
	DB.AutoMigrate(&model.Thread{})
	DB.AutoMigrate(&model.Torrent{})
	// Check if the group exists, if not, create it
	var count int64
	DB.Model(&model.UserGroup{}).Where("name = ?", group.Name).Count(&count)
	if count == 0 {
		DB.Create(&group)
	}

	return nil
}
