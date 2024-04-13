package service

import (
	"context"
	"fmt"
	"io"
	"s3-service/internal/entities"
	"strings"

	"github.com/google/uuid"
)

type Storage interface {
	CreateFile(ctx context.Context, fileInfo *entities.FileInfo, reader io.Reader) error
}

type Service struct {
	storage Storage
}

func New(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) CreateFile(ctx context.Context, file *entities.File) error {
	id, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("Service.CreateFile unable to create UUID: %w", err)
	}

	fmt.Printf("Service.CreateFile file ID: %s\n", file.ContentType)

	file.Name = fmt.Sprintf("%s.%s", id.String(), strings.Split(file.ContentType, "/")[1])

	if err := s.storage.CreateFile(ctx, &file.FileInfo, &file.Buffer); err != nil {
		return fmt.Errorf("Service.CreateFile unable to create file: %w", err)
	}

	return nil
}
