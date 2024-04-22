package service

import (
	"bytes"
	"context"
	"fmt"
	"image/jpeg"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/t1d333/smartlectures/internal/recognizer"
	"github.com/t1d333/smartlectures/pkg/logger"

	"github.com/karmdip-mi/go-fitz"
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

func (s *Service) ImportPdf(file []byte, ctx context.Context) (string, error) {
	reader := bytes.NewReader(file)
	doc, err := fitz.NewFromReader(reader)
	imgs := [][]byte{}
	if err != nil {
		return "", fmt.Errorf("failed to read pdf file: %w", err)
	}

	for n := 0; n < doc.NumPage(); n++ {
		buf := bytes.NewBuffer([]byte{})
		img, err := doc.Image(n)

		if err != nil {
			return "", fmt.Errorf("failed to convert %d page to image: %w", n, err)
		}

		if err = jpeg.Encode(buf, img, &jpeg.Options{Quality: jpeg.DefaultQuality}); err != nil {
			return "", fmt.Errorf("failed to encode jpeg page: %w", err)
		}

		imgs = append(imgs, buf.Bytes())
	}

	return s.RecognizeMixed(imgs, ctx)
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
