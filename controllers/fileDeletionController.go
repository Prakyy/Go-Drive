package controllers

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prakyy/Go-Drive/initializers"
	"github.com/prakyy/Go-Drive/models"
)

// DeleteFile handles file deletion
func DeleteFile(c *gin.Context) {
	// Get the file ID from the URL parameter
	fileID := c.Param("file_id")

	// Fetch file metadata from the database
	var fileMetadata models.File
	result := initializers.DB.First(&fileMetadata, "id = ?", fileID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Check if the logged-in user is the owner of the file
	user, _ := c.Get("user")
	//user := userModel.(models.User)

	if strconv.FormatUint(uint64(user.(models.User).ID), 10) != fileMetadata.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this file"})
		return
	}

	// Delete the file from storage (S3/local directory)
	err := os.Remove(fileMetadata.StoragePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete the file from server"})
		return
	}

	// Remove file metadata from the database
	initializers.DB.Unscoped().Delete(&fileMetadata)

	// Success response
	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}
