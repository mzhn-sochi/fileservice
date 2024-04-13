package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"io"
	"s3-service/internal/entities"
	"s3-service/pkg/utils"
	"strings"
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
	id, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("Service.CreateFile unable to create UUID: %w", err)
	}

	mime, ok := utils.GetContentTypeFromB64(file.Buffer.String())
	if !ok {
		return fmt.Errorf("unknown content-type")
	}

	file.ContentType = mime

	file.Name = fmt.Sprintf("%s.%s", id.String(), strings.Split(file.ContentType, "/")[1])

	if err := s.storage.CreateFile(ctx, &file.FileInfo, &file.Buffer); err != nil {
		return fmt.Errorf("Service.CreateFile unable to create file: %w", err)
	}

	return nil
}
