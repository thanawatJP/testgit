package handler

import (
	"authenservice/database"
	"authenservice/database/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateRoleHandler(c *gin.Context) {
	var newRole models.Role

	if err := c.ShouldBindJSON(&newRole); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := database.DB.Create(&newRole)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusCreated, newRole)
}

func GetAllRoleHandler(c *gin.Context) {
	var role []models.Role
	result := database.DB.Find(&role)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
	}
	c.JSON(http.StatusOK, role)
}
