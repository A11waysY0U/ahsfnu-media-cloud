package database

import (
	"fmt"
	"log"

	"ahsfnu-media-cloud/internal/config"
	"ahsfnu-media-cloud/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() {
	cfg := config.AppConfig.Database

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Shanghai",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")

	// 自动迁移所有模型
	err = DB.AutoMigrate(
		&models.User{},
		&models.InviteCode{},
		&models.Material{},
		&models.Tag{},
		&models.MaterialTag{},
		&models.WorkflowGroup{},
		&models.WorkflowMember{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migration completed successfully")
}

func GetDB() *gorm.DB {
	return DB
}
