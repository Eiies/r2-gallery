package controllers

import (
	"net/http"
	"r2-gallery/config"
	"r2-gallery/models"
	"r2-gallery/services"
	"r2-gallery/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=3"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Register handles user registration
func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.LogError("无法绑定 JSON", err)
		utils.SendError(c, http.StatusBadRequest, "无法绑定 JSON")
		return
	}

	// Check if user already exists
	var existingUser models.User
	result := config.DB.Where("email = ?", input.Email).Or("username = ?", input.Username).First(&existingUser)
	if result.RowsAffected > 0 {
		utils.SendError(c, http.StatusConflict, "此电子邮件或用户名的用户已存在")
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.LogError("无法对密码进行哈希处理", err)
		utils.SendError(c, http.StatusInternalServerError, "无法对密码进行哈希处理")
		return
	}

	// Create user
	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     "user",
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// Login handles user login
func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user
	var user models.User
	result := config.DB.Where("email = ?", input.Email).First(&user)
	if result.Error != nil {
		utils.SendError(c, http.StatusUnauthorized, "无效的电子邮件或密码")
		return
	}

	// Verify password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		utils.SendError(c, http.StatusUnauthorized, "无效的电子邮件或密码")
		return
	}

	// Generate token
	token, err := services.GenerateToken(user.ID, user.Role)
	if err != nil {
		utils.LogError("无法生成 Token", err)
		utils.SendError(c, http.StatusInternalServerError, "无法生成 Token")
		return
	}

	utils.SendSuccess(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}
