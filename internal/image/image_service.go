package image

import (
	"bytes"
	"time"
)

type service struct {
	repository *repository
	timeout    time.Duration
}

func NewService(r *repository) *service {
	return &service{
		repository: r,
		timeout:    time.Duration(10) * time.Second,
	}
}

func (s *service) UploadFileToMongoDB(req *ImageUploadRequest) (*ImageUploadResponse, error) {
	file, err := req.File.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()
	s.repository.UploadImage(file, req.File.Filename)
	res := &ImageUploadResponse{
		ImageName: req.File.Filename,
		Message:   "File uploaded",
	}
	return res, nil
}

func (s *service) DownloadFile(filename *string) (*bytes.Buffer, error) {
	return s.repository.DownloadImage(*filename)
}
