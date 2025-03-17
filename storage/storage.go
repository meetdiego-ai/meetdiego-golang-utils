package storage

import (
	"bytes"
	"context"
	"os"
	"strconv"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/redis/go-redis/v9"
)

// MinioStorage handles operations with MinIO/S3 compatible storage
type MinioStorage struct {
	Client     *minio.Client
	BucketName string
}

type RedisStorage struct {
	Client *redis.Client
}

// NewMinioStorage creates a new MinioStorage instance
func NewMinioStorage(bucketName string) *MinioStorage {
	// Initialize minio client
	endpoint := os.Getenv("R2_ENDPOINT")
	if endpoint == "" {
		panic("R2_ENDPOINT is not set")
	}

	accessKeyID := os.Getenv("R2_ACCESS_KEY_ID")
	if accessKeyID == "" {
		panic("R2_ACCESS_KEY_ID is not set")
	}

	secretAccessKey := os.Getenv("R2_SECRET_ACCESS_KEY")
	if secretAccessKey == "" {
		panic("R2_SECRET_ACCESS_KEY is not set")
	}

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

func NewRedisClient() (*redis.Client, error) {
	if os.Getenv("REDIS_ADDR") == "" {
		panic("REDIS_ADDR is not set")

	}

	if os.Getenv("REDIS_PASSWORD") == "" {
		panic("REDIS_PASSWORD is not set")
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB: func() int {
			db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
			if err != nil {
				return 0 // Default to DB 0 if conversion fails
			}
			return db
		}(),
	})
	return redisClient, nil
}

func GetRedisClient() *redis.Client {
	redisClient, err := NewRedisClient()
	if err != nil {
		panic(err)
	}
	return redisClient
}
