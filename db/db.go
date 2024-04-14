package db

import (
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

type Feature struct {
	gorm.Model
	Description string
}

func SetupDatabase() *gorm.DB {
	dsn := "host=localhost user=daze dbname=bannerservice sslmode=disable password=devps123"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&Banner{}, &BannerVersion{}, &Tag{}, &Feature{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	return db
}
