package controllers

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prakyy/Go-Drive/initializers"
	"github.com/prakyy/Go-Drive/models"
)

func GetFile(c *gin.Context) {
	// Get the file ID from the URL parameter
	fileID := c.Param("file_id")

	// Fetch file metadata from the database
	var fileMetadata models.File
	result := initializers.DB.First(&fileMetadata, "id = ?", fileID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Serve the file
	filePath := fileMetadata.StoragePath
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found on server"})
		return
	}

	// Check if the file is public
	if fileMetadata.IsPublic {
		// Serve the file to public
		c.File(filePath)
		return
	}

	// Check file visibility
	user, exists := c.Get("user")
	if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }
    // Type assertion to extract UserID from models.User
    userModel, ok := user.(models.User)
    if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user data"})
        return
    }
	
	useriduint := strconv.FormatUint(uint64(userModel.ID), 10)
    if useriduint != fileMetadata.UserID {
        c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
        return
    }

	c.File(filePath)
}