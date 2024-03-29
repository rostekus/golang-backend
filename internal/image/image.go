package image

import (
	"mime/multipart"
)

type Image struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
}

type ImageMessage struct {
	ImageID string `json:"image_id"`
	UserID  string `json:"user_id"`
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
