package storage

import (
	"bytes"
	"context"

	"github.com/minio/minio-go/v7"
)

// MinioStorage handles operations with MinIO/S3 compatible storage
type MinioStorage struct {
	Client     *minio.Client
	BucketName string
}

// NewMinioStorage creates a new MinioStorage instance
func NewMinioStorage(client *minio.Client, bucketName string) *MinioStorage {
	return &MinioStorage{
		Client:     client,
		BucketName: bucketName,
	}
}

// SaveContent saves the provided content to the specified object path
func (s *MinioStorage) SaveContent(objectName, content string) error {
	_, err := s.Client.PutObject(
		context.Background(),
		s.BucketName,
		objectName,
		bytes.NewReader([]byte(content)),
		int64(len(content)),
		minio.PutObjectOptions{
			ContentType: "text/plain",
		},
	)

	return err
}
