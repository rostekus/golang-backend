package image

import (
	"mime/multipart"
)

type ImageUploadRequest struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type ImageUploadResponse struct {
	ImageName string `json:"image_name"`
	Message   string `json:"message"`
}
type ImageDownloadRequest struct {
	ImageName string `json:"image_name"`
}
type ImageDownloadResponse struct {
}
