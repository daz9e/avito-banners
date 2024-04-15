package handlers

import (
	"avito-banners/internal/config"
	"avito-banners/internal/db"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func GetUserBanner(c *gin.Context) {
	tagID := c.Query("tag_id")
	featureID := c.Query("feature_id")
	useLastRevision := c.Query("use_last_revision") == "true"
	role, exists := c.Get("role")

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user role not found"})
		return
	}

	if tagID == "" || featureID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tag_id and feature_id are required"})
		return
	}
	cacheKey := fmt.Sprintf("banner:%s:%s", featureID, tagID)
	if !useLastRevision {
		val, err := config.Redis.Get(config.Ctx, cacheKey).Result()
		if err == nil {
			var banners []db.Banner
			if err := json.Unmarshal([]byte(val), &banners); err == nil {
				c.JSON(http.StatusOK, banners)
				return
			}
		}
	}

	var banners []db.Banner
	query := config.Database.Joins("JOIN banner_tags on banner_tags.banner_id = banners.id").
		Where("banner_tags.tag_id = ? AND banners.feature_id = ?", tagID, featureID)

	if role != "admin" {
		query = query.Where("banners.is_active = true")
	}

	result := query.Find(&banners)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if len(banners) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no banners found"})
		return
	}

	data, _ := json.Marshal(banners)
	config.Redis.Set(config.Ctx, cacheKey, data, 5*time.Minute)

	c.JSON(http.StatusOK, banners)
}

func GetBanners(c *gin.Context) {
	featureID, _ := strconv.Atoi(c.Query("feature_id"))
	tagID, _ := strconv.Atoi(c.Query("tag_id"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	query := config.Database.Model(&db.Banner{})

	if featureID > 0 {
		query = query.Where("feature_id = ?", featureID)
	}
	if tagID > 0 {
		query = query.Joins("JOIN banner_tags on banners.id = banner_tags.banner_id").Where("banner_tags.tag_id = ?", tagID)
	}

	var banners []db.Banner
	if err := query.Offset(offset).Limit(limit).Find(&banners).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, banners)
}

func CreateBanner(c *gin.Context) {
	var input struct {
		FeatureID uint   `json:"feature_id"`
		TagIDs    []uint `json:"tag_ids"`
		Content   string `json:"content"`
		IsActive  bool   `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Bind error": err.Error()})
		return
	}

	banner := db.Banner{
		FeatureID: input.FeatureID,
		Content:   input.Content,
		IsActive:  input.IsActive,
	}

	for _, tagID := range input.TagIDs {
		var tag db.Tag
		if err := config.Database.First(&tag, tagID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("tag with ID %d not found", tagID)})
			return
		}
		banner.Tags = append(banner.Tags, &tag)
	}

	if result := config.Database.Create(&banner); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, banner)
}

func UpdateBanner(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var banner db.Banner
	if err := config.Database.First(&banner, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Banner not found"})
		return
	}

	var input struct {
		Content  *string `json:"content"`
		IsActive *bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Создание новой версии баннера перед изменениями
	newVersion := db.BannerVersion{
		BannerID:  banner.ID,
		FeatureID: banner.FeatureID,
		Content:   banner.Content,
		IsActive:  banner.IsActive,
	}

	if err := config.Database.Create(&newVersion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create banner version: " + err.Error()})
		return
	}

	if input.Content != nil {
		banner.Content = *input.Content
	}
	if input.IsActive != nil {
		banner.IsActive = *input.IsActive
	}

	if err := config.Database.Save(&banner).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update banner: " + err.Error()})
		return
	}

	if err := config.Database.Exec("DELETE FROM banner_versions WHERE id NOT IN (SELECT id FROM banner_versions WHERE banner_id = ? ORDER BY created_at DESC LIMIT 3)", banner.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete old versions: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, banner)
}

func DeleteBanner(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid banner ID"})
		return
	}

	result := config.Database.Where("id = ?", id).Delete(&db.Banner{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "banner not found"})
		return
	}
	c.Status(http.StatusNoContent)
}

func GetBannerVersions(c *gin.Context) {
	bannerID, err := strconv.Atoi(c.Param("banner_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid banner ID"})
		return
	}

	var versions []db.BannerVersion
	if err := config.Database.Where("banner_id = ?", bannerID).Order("created_at desc").Find(&versions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(versions) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No versions found for the banner"})
		return
	}

	c.JSON(http.StatusOK, versions)
}

func RestoreBannerVersion(c *gin.Context) {
	versionID, err := strconv.Atoi(c.Param("version_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid version ID"})
		return
	}

	var version db.BannerVersion
	if err := config.Database.First(&version, versionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Version not found"})
		return
	}

	var banner db.Banner
	if err := config.Database.First(&banner, version.BannerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Banner not found"})
		return
	}

	banner.FeatureID = version.FeatureID
	banner.Content = version.Content
	banner.IsActive = version.IsActive
	if err := config.Database.Save(&banner).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Banner has been restored to the selected version", "banner": banner})
}
