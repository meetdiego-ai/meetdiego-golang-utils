package storage

import (
	"bytes"
	"context"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioStorage handles operations with MinIO/S3 compatible storage
type MinioStorage struct {
	Client     *minio.Client
	BucketName string
}

// NewMinioStorage creates a new MinioStorage instance
func NewMinioStorage(bucketName string) *MinioStorage {
	// Initialize minio client
	endpoint := os.Getenv("R2_ENDPOINT")
	accessKeyID := os.Getenv("R2_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("R2_SECRET_ACCESS_KEY")
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		panic(err)
	}

	return &MinioStorage{
		Client:     client,
		BucketName: bucketName,
	}
}

// SaveContent saves the provided content to the specified object path
func (s *MinioStorage) SaveContent(objectName, content string, contentType string) error {
	if contentType == "" {
		contentType = "text/plain"
	}

	_, err := s.Client.PutObject(
		context.Background(),
		s.BucketName,
		objectName,
		bytes.NewReader([]byte(content)),
		int64(len(content)),
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)

	return err
}
