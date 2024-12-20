package handler

import (
	"authenservice/database"
	"authenservice/database/models"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
)

func RegisterUser(newUser models.UserAuth) (*models.UserAuth, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}
	newUser.Password = string(hashedPassword)

	var role models.Role
	if err := database.DB.First(&role, newUser.RoleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}

	newUser.Role = role

	result := database.DB.Create(&newUser)
	if result.Error != nil {
		return nil, result.Error
	}

	return &newUser, nil
}

func CreateUserHandler(c *gin.Context) {
	var newUser models.UserAuth

	// Bind JSON input to newUser
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Register user
	createdUser, err := RegisterUser(newUser)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, gorm.ErrRecordNotFound) {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}

func GetAllUserHandler(c *gin.Context) {
	var users []models.UserAuth
	result := database.DB.Preload("Role").Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func GetOneUserHandler(c *gin.Context) {
	id := c.Param("id")
	var user models.UserAuth

	result := database.DB.Preload("Role").First(&user, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
