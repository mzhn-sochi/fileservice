package storage

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"log"
	"os"
	"s3-service/internal/entities"
)

type MinioStorage struct {
	client *minio.Client
	bucket string
}

func New(sslEnabled bool) (*MinioStorage, error) {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")

	log.Printf("Endpoint: %s", endpoint)
	log.Printf("Access Key: %s", accessKey)
	log.Printf("Secret Key: %s", secretKey)

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: sslEnabled,
	})

	if err != nil {
		return nil, err
	}

	return &MinioStorage{
		client: client,
		bucket: os.Getenv("MINIO_BUCKET"),
	}, nil
}

func (s *MinioStorage) checkBucket(ctx context.Context) error {
	exists, err := s.client.BucketExists(ctx, s.bucket)
	if err != nil {
		return fmt.Errorf("MinioStorage cannot check bucket existance: %w", err)
	}

	if !exists {
		if err := s.client.MakeBucket(ctx, s.bucket, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("MinioStorage cannot create bucket: %w", err)
		}
	}

	return nil
}

func (s *MinioStorage) CreateFile(ctx context.Context, fileInfo *entities.FileInfo, reader io.Reader) error {
	if err := s.checkBucket(ctx); err != nil {
		return fmt.Errorf("MinioStorage.CreateFile: %w", err)
	}

	// TODO Можно запихнуть метадату
	_, err := s.client.PutObject(ctx, s.bucket, fileInfo.Name, reader, fileInfo.Size, minio.PutObjectOptions{
		ContentType: fileInfo.ContentType,
	})
	if err != nil {
		return fmt.Errorf("MinioStorage.CreateFile unable to put object: %w", err)
	}

	return nil
}
