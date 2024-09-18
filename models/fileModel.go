package models

import (
	"time"

	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	ID          string    `gorm:"primaryKey" json:"id"`
	UserID      string    `json:"user_id"`
	FileName    string    `json:"file_name"`
	FileSize    int64     `json:"file_size"`
	FileType    string    `json:"file_type"`
	StoragePath string    `json:"storage_path"`
	UploadDate  time.Time `json:"upload_date"`
	IsPublic    bool      `json:"is_public"`
}