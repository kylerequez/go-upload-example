package models

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID        uuid.UUID
	Name      string
	Size      int64
	Type      string
	CreatedAt time.Time
}

func NewFile(fileName string, fileSize int64, fileType string, createdAt time.Time) *File {
	return &File{
		Name:      fileName,
		Size:      fileSize,
		Type:      fileType,
		CreatedAt: createdAt,
	}
}
