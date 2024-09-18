package controllers

import (
	"fmt"
	"strconv"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prakyy/Go-Drive/initializers"
	"github.com/prakyy/Go-Drive/models"
)

const uploadDir = "./s3bucket"

func UploadFile(c *gin.Context) {
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
