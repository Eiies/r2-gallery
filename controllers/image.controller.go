package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"r2-gallery/config"
	"r2-gallery/models"
	"r2-gallery/services"
	"r2-gallery/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type ImageResponse struct {
	ID        uint      `json:"id"`
	FileName  string    `json:"file_name"`
	URL       string    `json:"url"`
	Title     string    `json:"title"`
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

// UploadImage handles image upload
func UploadImage(c *gin.Context) {
	// Get user from context
	claims, exists := c.Get("user")
	if !exists {
		utils.LogError("Unauthorized", nil)
		utils.SendError(c, http.StatusUnauthorized, "没有权限")
		return
	}

	userClaims := claims.(jwt.MapClaims)
	userID, _ := strconv.ParseUint(fmt.Sprintf("%v", userClaims["user_id"]), 10, 32)

	// Get file from form
	file, err := c.FormFile("image")
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "未提供图片")
		return
	}

	// Get title from form
	title := c.PostForm("title")
	if title == "" {
		utils.SendError(c, http.StatusBadRequest, "标题为必填项")
		return
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%d_%s%s", time.Now().Unix(), title, ext)

	// Open file
	src, err := file.Open()
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "无法打开文件")
		return
	}
	defer src.Close()

	// Upload to R2
	url, err := services.UploadToR2(src, fileName)
	if err != nil {
		utils.LogError("上传文件失败", err)
		utils.SendError(c, http.StatusBadRequest, "上传文件失败")
		return
	}

	// Save to database
	image := models.Image{
		FileName:  fileName,
		URL:       url,
		Title:     title,
		UserID:    uint(userID),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := config.DB.Create(&image).Error; err != nil {
		// Try to delete from R2 if database save fails
		_ = services.DeleteFromR2(fileName)
		utils.LogError("无法保存图像", err)
		utils.SendError(c, http.StatusInternalServerError, "无法保存图像")
		return
	}

	utils.SendSuccess(c, gin.H{
		"image": ImageResponse{
			ID:        image.ID,
			FileName:  image.FileName,
			URL:       image.URL,
			Title:     image.Title,
			UserID:    image.UserID,
			CreatedAt: image.CreatedAt,
		},
	})
}

// ListImages returns all images for the current user
func ListImages(c *gin.Context) {
	// Get user from context
	claims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userClaims := claims.(jwt.MapClaims)
	userID, _ := strconv.ParseUint(fmt.Sprintf("%v", userClaims["user_id"]), 10, 32)

	// Get images from database
	var images []models.Image
	result := config.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&images)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch images"})
		return
	}

	// Format response
	var response []ImageResponse
	for _, image := range images {
		// Get username
		var user models.User
		config.DB.Select("username").First(&user, image.UserID)

		response = append(response, ImageResponse{
			ID:        image.ID,
			FileName:  image.FileName,
			URL:       image.URL,
			Title:     image.Title,
			UserID:    image.UserID,
			Username:  user.Username,
			CreatedAt: image.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"images": response})
}

// DeleteImage deletes an image
func DeleteImage(c *gin.Context) {
	// Get user from context
	claims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userClaims := claims.(jwt.MapClaims)
	userID, _ := strconv.ParseUint(fmt.Sprintf("%v", userClaims["user_id"]), 10, 32)

	// Get image ID from URL
	imageID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image ID"})
		return
	}

	// Get image from database
	var image models.Image
	result := config.DB.First(&image, imageID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	// Check if user owns the image
	if image.UserID != uint(userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this image"})
		return
	}

	// Delete from R2
	err = services.DeleteFromR2(image.FileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image from storage"})
		return
	}

	// Delete from database
	result = config.DB.Delete(&image)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image from database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image deleted successfully"})
}
