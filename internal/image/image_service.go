package image

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"rostekus/golang-backend/pkg/db"
	"rostekus/golang-backend/pkg/rabbitmq"
	"time"

	"github.com/google/uuid"
)

type storage interface {
	UploadFile(file multipart.File, filename string) error
	DownloadFile(filename string) (*bytes.Buffer, error)
}

type publisher interface {
	Publish(message []byte) error
}

type service struct {
	imageStorage storage
	db           db.SQLDatabase
	publisher    publisher
	timeout      time.Duration
}

func NewService(r *repository, p *rabbitmq.RabbitMQ, db db.SQLDatabase) *service {
	return &service{
		imageStorage: r,
		timeout:      time.Duration(10) * time.Second,
		publisher:    p,
		db:           db,
	}
}

func (s *service) UploadFileToMongoDB(ctx context.Context, req *ImageUploadRequest) (*ImageUploadResponse, error) {

	file, err := req.File.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()
	imageID := uuid.New()
	s.imageStorage.UploadFile(file, imageID.String())
	res := &ImageUploadResponse{
		ImageName: req.File.Filename,
		Message:   "File uploaded",
		ImageID:   imageID.String(),
	}
	return res, nil
}

func (s *service) PublishMessage(imageMessage ImageMessage) error {
	body, err := json.Marshal(imageMessage)
	if err != nil {
		return err
	}
	return s.publisher.Publish(body)
}

func (s *service) InsertImageDataToDB(ctx context.Context, req *ImageUploadRequest, userID string, imageID string) error {
	query := "INSERT INTO images(id, user_id,filename) VALUES ($1, $2, $3)"
	_, err := s.db.ExecContext(ctx, query, imageID, userID, req.File.Filename)
	return err
}

func (s *service) DownloadFile(ctx context.Context, filename string, userID string) (*ImageDownloadResponse, *bytes.Buffer, error) {

	query := "SELECT id, filename FROM images WHERE id = $1 AND user_id = $2"
	row := s.db.QueryRowContext(ctx, query, filename, userID)

	var id string
	var foundFilename string
	err := row.Scan(&id, &foundFilename)
	if err == sql.ErrNoRows {
		// No image found with the given filename and user ID.
		return nil, nil, fmt.Errorf("no image found with filename %s and user ID %s", filename, userID)
	} else if err != nil {
		// Error occurred while querying the database.
		return nil, nil, err
	}
	buf, err := s.imageStorage.DownloadFile(id)
	resp := &ImageDownloadResponse{
		ImageName: filename,
		Message:   "File downloaded",
	}
	return resp, buf, err
}

func (s *service) GetImagesForUser(ctx context.Context, userID string) ([]*Image, error) {
	var images []*Image
	query := "SELECT filename, id FROM images WHERE user_id = $1"
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var image Image
		if err := rows.Scan(&image.Filename, &image.ID); err != nil {
			return nil, err
		}
		images = append(images, &image)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return images, nil
}
