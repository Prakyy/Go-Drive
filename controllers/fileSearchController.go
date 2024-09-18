package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prakyy/Go-Drive/initializers"
	"github.com/prakyy/Go-Drive/models"
)

func SearchFile(c *gin.Context) {
	// Get the file name from the query parameter
	fileName := c.Query("file")

	// Get authenticated user
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	userModel, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user data"})
		return
	}

	// Search the file in the database
	var fileMetadata models.File
	result := initializers.DB.Where("user_id = ? AND file_name = ?", strconv.FormatUint(uint64(userModel.ID), 10), fileName).First(&fileMetadata)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Return the file details
	c.JSON(http.StatusOK, gin.H{
		"file_id":     fileMetadata.ID,
		"filename":    fileMetadata.FileName,
		"upload_date": fileMetadata.CreatedAt,
	})
}