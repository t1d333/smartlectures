package elasticsearch

import "github.com/t1d333/smartlectures/internal/models"

type NoteSearchResponse struct {
	Took int64
	Hits struct {
		Total struct {
			Value int64
		}
		Hits []*NoteSearchHit
	}
}

type NoteSearchHit struct {
	Score     float64          `json:"_score"`
	Index     string           `json:"_index"`
	Type      string           `json:"_type"`
	Version   int64            `json:"_version,omitempty"`
	Highlight NoteHighlight    `json:"highlight"`
	Source    NoteSearchSource `json:"_source"`
}

type NoteHighlight struct {
	Name []string `json:"name"`
	Body []string `json:"body"`
}

type NoteSearchSource struct {
	Id   int    `json:"noteId"`
	Name string `json:"name"`
}

type SnippetSearchResponse struct {
	Took int64
	Hits struct {
		Total struct {
			Value int64
		}
		Hits []*SnippetSearchHit
	}
}

type SnippetSearchHit struct {
	Score     float64        `json:"_score"`
	Index     string         `json:"_index"`
	Type      string         `json:"_type"`
	Version   int64          `json:"_version,omitempty"`
	Source    models.Snippet `json:"_source"`
}
