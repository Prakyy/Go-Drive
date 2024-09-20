package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prakyy/Go-Drive/initializers"
	"github.com/prakyy/Go-Drive/models"
)


func ReceiveFile(c *gin.Context) {
	var uploadDir = os.Getenv("StorageBucket")
	
	// Get the user from the context (set by your auth middleware)
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Type assertion to models.User
	userModel, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user data"})
		return
	}

	// Parse the multipart form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to get file from form"})
		return
	}
	defer file.Close()

	// Parse public flag
	isPublic := c.PostForm("public") == "true"
	
	// Generate a unique filename
	filename := uuid.New().String() + filepath.Ext(header.Filename)

	// Convert user ID to string
	userID := strconv.FormatUint(uint64(userModel.ID), 10)
	
	// Ensure the upload directory exists
	userUploadDir := filepath.Join(uploadDir, userID)
	fmt.Println("User Upload Dir:", userUploadDir)  // Ensure path looks correct
	
	err = os.MkdirAll(userUploadDir, 0755)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	// Create the file
	filePath := filepath.Join(userUploadDir, filename)
	dst, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the destination file
	_, err = io.Copy(dst, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Get file size
	fileInfo, err := dst.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file info"})
		return
	}

	// Store metadata in the database
	fileMetadata := models.File{
		ID:          uuid.New().String(),
		UserID:      userID, // Ensure this matches the type in your database schema
		FileName:    header.Filename,
		FileSize:    fileInfo.Size(),
		FileType:    header.Header.Get("Content-Type"),
		StoragePath: filePath,
		UploadDate:  time.Now(),
		IsPublic:    isPublic,
	}

	result := initializers.DB.Create(&fileMetadata)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store file metadata"})
		return
	}

	// Return response
	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"file_id": fileMetadata.ID,
		"path":    fileMetadata.StoragePath,
	})
}

func ServeFile(c *gin.Context) {
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
