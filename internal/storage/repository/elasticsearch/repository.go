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

func (r *Repository) SearchNote(
	ctx context.Context,
	query string,
) ([]models.NoteSearchItem, error) {
	var buf bytes.Buffer

	// body := map[string]interface{}{
	// 	"query": map[string]interface{}{
	// 		"multi_match": map[string]interface{}{
	// 			"query":     query,
	// 			"fields":    []string{"name", "body"},
	// 			"fuzziness": "AUTO",
	// 		},
	// 	},
	// 	"highlight": map[string]interface{}{
	// 		"fields": map[string]interface{}{
	// 			"name": map[string]interface{}{
	// 				"fragment_size": 50,
	// 			},
	// 			"body": map[string]interface{}{
	// 				"fragment_size": 100,
	// 			},
	// 		},
	// 	},
	// }

	body := map[string]interface{}{
    "query": map[string]interface{}{
        "bool": map[string]interface{}{
            "should": []interface{}{
                map[string]interface{}{
                    "wildcard": map[string]interface{}{
                        "field_name": "*query*",
                    },
                },
                map[string]interface{}{
                    "multi_match": map[string]interface{}{
                        "query":     "текст запроса",
                        "fields":    []string{"name", "body"},
                        "fuzziness": "AUTO",
                    },
                },
            },
        },
    },
    "highlight": map[string]interface{}{
        "fields": map[string]interface{}{
            "name": map[string]interface{}{
                "fragment_size": 50,
            },
            "body": map[string]interface{}{
                "fragment_size": 100,
            },
        },
    },
}

	if err := jsoniter.NewEncoder(&buf).Encode(&body); err != nil {
		return []models.NoteSearchItem{}, fmt.Errorf("failed to create search request: %w", err)
	}

	res, err := r.client.Search(
		r.client.Search.WithBody(strings.NewReader(buf.String())),
		r.client.Search.WithSource("noteId", "name"),
	)
	if err != nil {
		return []models.NoteSearchItem{}, fmt.Errorf("failed to make search request: %w", err)
	}

	defer res.Body.Close()

	if res.IsError() {
		return []models.NoteSearchItem{}, fmt.Errorf("failed to search note")
	}

	var response NoteSearchResponse

	if err := jsoniter.NewDecoder(res.Body).Decode(&response); err != nil {
		return []models.NoteSearchItem{}, fmt.Errorf("failed to decode search result: %w", err)
	}

	result := []models.NoteSearchItem{}

	for _, hit := range response.Hits.Hits {
		result = append(result, models.NoteSearchItem{
			NoteID:        hit.Source.Id,
			Name:          hit.Source.Name,
			BodyHighlight: hit.Highlight.Body,
			NameHighlight: hit.Highlight.Name,
		})
	}

	return result, nil
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
