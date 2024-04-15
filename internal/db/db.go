package db

import (
	"avito-banners/internal/config"
	"fmt"
	"github.com/charmbracelet/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Banner struct {
	gorm.Model
	FeatureID uint   `gorm:"not null"`
	Tags      []*Tag `gorm:"many2many:banner_tags;"`
	Content   string `gorm:"type:text;not null"`
	IsActive  bool   `gorm:"not null;default:true"`
}

type BannerVersion struct {
	gorm.Model
	BannerID  uint   `gorm:"index;not null;onDelete:cascade"`
	FeatureID uint   `gorm:"not null"`
	Content   string `gorm:"type:text;not null"`
	IsActive  bool   `gorm:"not null"`
}

type Tag struct {
	gorm.Model
	Description string
	Banners     []*Banner `gorm:"many2many:banner_tags;"`
}

func SetupDatabase() {

	cfg := config.Cfg

	dsn := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=%s password=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBName, cfg.DBSSLMode, cfg.DBPassword)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&Banner{}, &BannerVersion{}, &Tag{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	config.Database = db

	err = seedTags()
	if err != nil {
		log.Fatalf("error seeding tags: %v", err)
	}
}

func seedTags() error {
	tags := []Tag{
		{Description: "Technology"},
		{Description: "Health"},
		{Description: "Finance"},
		{Description: "Education"},
	}

	for _, tag := range tags {
		result := config.Database.FirstOrCreate(&tag, Tag{Description: tag.Description})
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}
