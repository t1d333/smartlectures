package elasticsearch

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"

	elastic "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	jsoniter "github.com/json-iterator/go"
	"github.com/t1d333/smartlectures/internal/models"
	"github.com/t1d333/smartlectures/internal/storage/repository"
	"github.com/t1d333/smartlectures/pkg/logger"
)

type Repository struct {
	client *elastic.Client
	logger logger.Logger
}

func (*Repository) SearchDir(ctx context.Context, query string) ([]int, error) {
	panic("unimplemented")
}

// SearchNoteByBody implements repository.Repository.
func (*Repository) SearchNoteByBody(ctx context.Context, query string) ([]int, error) {
	panic("unimplemented")
}

// SearchNoteByName implements repository.Repository.
func (r *Repository) SearchNoteByName(ctx context.Context, query string) ([]int, error) {
	body := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"name": query,
			},
		},
	}

	var rawBody bytes.Buffer

	if err := jsoniter.NewEncoder(&rawBody).Encode(body); err != nil {
		return []int{}, fmt.Errorf("failed to create search query: %w", err)
	}

	req := esapi.SearchRequest{
		Index: []string{"your_index"},
		Body:  strings.NewReader(rawBody.String()),
	}

	// Выполнение запроса
	res, err := req.Do(context.Background(), r.client)
	if err != nil {
		return []int{}, fmt.Errorf("failed to make search request: %w", err)
	}

	defer res.Body.Close()

	var tmp map[string]interface{}

	if err := jsoniter.NewDecoder(res.Body).Decode(&tmp); err != nil {
	}

	fmt.Println(r)

	return []int{}, nil
}

func NewRepository(logger logger.Logger) (repository.Repository, error) {
	client, err := elastic.NewDefaultClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create elastic client: %w", err)
	}

	_, err = client.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to elasticsearch: %w", err)
	}

	return &Repository{
		client: client,
		logger: logger,
	}, nil
}

func (r *Repository) GetNote(ctx context.Context, id int) (models.Note, error) {
	req := esapi.GetRequest{
		Index:      "notes",
		DocumentID: strconv.Itoa(id),
		Source:     []string{"name", "body"},
		FilterPath: []string{"_source"},
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		return models.Note{}, fmt.Errorf("failed to get note in repository: %w", err)
	}

	defer res.Body.Close()

	var data struct {
		Source models.Note `json:"_source"`
	}

	var body bytes.Buffer
	body.ReadFrom(res.Body)

	if err != nil {
		return models.Note{}, fmt.Errorf("failed to read note: %w", err)
	}

	if jsoniter.Unmarshal(body.Bytes(), &data); err != nil {
		return models.Note{}, fmt.Errorf("failed to unmarshal note: %w", err)
	}

	return data.Source, nil
}

func (r *Repository) CreateNote(ctx context.Context, note models.Note) error {
	raw, err := jsoniter.Marshal(note)
	if err != nil {
		return fmt.Errorf("failed to marshal note struct: %w", err)
	}

	req := esapi.CreateRequest{
		Index:      "notes",
		Body:       bytes.NewReader(raw),
		DocumentID: strconv.Itoa(note.NoteId),
		Refresh:    "true",
	}

	if _, err = req.Do(ctx, r.client); err != nil {
		return fmt.Errorf("failed to store note in repository: %w", err)
	}

	return nil
}

func (r *Repository) UpdateNote(ctx context.Context, note models.Note) error {
	r.logger.Info(note)
	body, err := jsoniter.Marshal(note)
	if err != nil {
		return fmt.Errorf("failed to marshal note: %w", err)
	}

	req := esapi.UpdateRequest{
		Index:      "notes",
		DocumentID: strconv.Itoa(note.NoteId),
		Body:       bytes.NewReader([]byte(fmt.Sprintf(`{"doc":%s}`, body))),
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		return fmt.Errorf("failed to make update note request: %w", err)
	}

	if res.IsError() {
		return fmt.Errorf("failed update note: " + res.Status())
	}

	return nil
}

func (r *Repository) DeleteNote(ctx context.Context, id int) error {
	r.logger.Info(id)
	req := esapi.DeleteRequest{
		Index:      "notes",
		DocumentID: strconv.Itoa(id),
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		return fmt.Errorf("failed to make delete request: %w", err)
	}

	if res.IsError() {
		return fmt.Errorf("failed to update note data in repository: %w", err)
	}

	return nil
}

func (r *Repository) Search(ctx context.Context, query string) error {
	panic("unimplemented")
}
