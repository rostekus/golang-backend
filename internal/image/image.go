package image

import (
	"mime/multipart"

	"github.com/google/uuid"
)

type Image struct {
	ID       uuid.UUID
	Filename string
}

type ImageUploadRequest struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type ImageUploadResponse struct {
	ImageName string `json:"image_name"`
	ImageID   string `json:"image_id"`
	Message   string `json:"message"`
}
type ImageDownloadRequest struct {
	ImageName string `form:"image_file"`
}
type ImageDownloadResponse struct {
	ImageName string `json:"image_name"`
	Message   string `json:"message"`
}
