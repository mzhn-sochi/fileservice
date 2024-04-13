package handler

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"s3-service/api/s3"
	"s3-service/internal/entities"
)

type Service interface {
	CreateFile(ctx context.Context, file *entities.File) error
}

type Handler struct {
	s3.UnimplementedS3Server

	service Service
}

func New(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h Handler) Upload(stream s3.S3_UploadServer) error {
	file := &entities.File{}

	writer := bufio.NewWriter(&file.Buffer)

	for {
		object, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Printf("%s", err.Error())
			return status.Error(codes.Internal, fmt.Sprintf("%s", err.Error()))
		}

		log.Printf("Metadata: %+v", object.Meta)

		chunk := object.Image.Chunk

		file.ContentType = object.Meta.ContentType
		file.Size += int64(len(chunk))
		if _, err := writer.Write(chunk); err != nil {
			fmt.Printf("%s", err.Error())
			return err
		}
	}

	err := h.service.CreateFile(stream.Context(), file)
	if err != nil {
		return err
	}

	return stream.SendAndClose(&s3.ObjectInfo{
		Name: file.Name,
	})
}
