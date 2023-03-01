package image

import (
	"bytes"
	"mime/multipart"
	"rostekus/golang-backend/db"
)

type repository struct {
	db           *db.MongoDatabase
	databaseName string
}

func NewRepository(db *db.MongoDatabase, databaseName string) *repository {
	return &repository{
		db:           db,
		databaseName: databaseName,
	}
}

func (r *repository) UploadFile(file multipart.File, filename string) error {
	return r.db.UploadFile(file, filename, r.databaseName)
}

func (r *repository) DownloadFile(filename string) (*bytes.Buffer, error) {
	return r.db.DownloadFile(filename, r.databaseName)
}
