package main

import (
	"github.com/t1d333/smartlectures/internal/storage/repository/elasticsearch"
	"github.com/t1d333/smartlectures/pkg/logger"
)

func main() {
	logger := logger.NewLogger()
	_, err := elasticsearch.NewRepository(logger)
	if err != nil {
		logger.Fatalf("failed to create elasticsearch repository: %s", err)
	}
}
