package service

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/t1d333/smartlectures/internal/recognizer"
	"github.com/t1d333/smartlectures/pkg/logger"
)

type Service struct {
	logger logger.Logger
	client RecognizerClient
}

func (s *Service) RecognizeFormula(img []byte, ctx context.Context) (string, error) {
	result, err := s.client.RecognizeFormula(ctx, &ImageToRecognize{Data: img})
	if err != nil {
		return "", fmt.Errorf("failed to recognize formula in service: %w", err)
	}

	return result.GetResult(), nil
}

func (s *Service) RecognizeMixed(imgs [][]byte, ctx context.Context) (string, error) {
	result, err := s.client.RecognizeMixed(ctx, &ImagesArrToRecognize{Data: imgs})
	if err != nil {
		return "", fmt.Errorf("failed to recognize mixed text in service: %w", err)
	}

	return result.GetResult(), nil
}

func (s *Service) RecognizeText(imgs [][]byte, ctx context.Context) (string, error) {
	result, err := s.client.RecognizeText(ctx, &ImagesArrToRecognize{Data: imgs})
	if err != nil {
		return "", fmt.Errorf("failed to recognize text in service: %w", err)
	}

	return result.GetResult(), nil
}

func NewService(logger logger.Logger, address string) recognizer.Service {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("did not connect: %v", err)
	}

	client := NewRecognizerClient(conn)

	return &Service{
		logger: logger,
		client: client,
	}
}
