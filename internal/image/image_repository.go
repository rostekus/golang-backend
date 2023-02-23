package image

import (
	"bytes"
	"mime/multipart"
)

type database interface {
	UploadFile(fileReader multipart.File, filename string, databaseName string) error
	DownloadFile(string, string) (*bytes.Buffer, error)
}

type repository struct {
	db           database
	databaseName string
}

func NewRepository(db database, databaseName string) *repository {
	return &repository{
		db:           db,
		databaseName: databaseName,
	}
}

func (r *repository) UploadImage(file multipart.File, filename string) error {
	return r.db.UploadFile(file, filename, r.databaseName)
}

func (r *repository) DownloadImage(filename string) (*bytes.Buffer, error) {
	return r.db.DownloadFile(filename, r.databaseName)
}
