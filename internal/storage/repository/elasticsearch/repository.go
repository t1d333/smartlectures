package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	elastic "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/t1d333/smartlectures/internal/storage/repository"
	"github.com/t1d333/smartlectures/pkg/logger"
)

type Repository struct {
	client *elastic.Client
	logger logger.Logger
}

func NewRepository(logger logger.Logger) (repository.Repository, error) {
	client, err := elastic.NewDefaultClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create elastic client: %w", err)
	}

	data, err := json.Marshal(struct {
		Title string `json:"title"`
	}{Title: "test1234"})

	req := esapi.IndexRequest{
		Index:      "test",
		DocumentID: strconv.Itoa(10),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), client)
	
	if err != nil {
		logger.Error("Error getting response: %s", err)
	}
	
	logger.Info(res)
	defer res.Body.Close()
	

	return Repository{
		client: client,
		logger: logger,
	}, nil
}
