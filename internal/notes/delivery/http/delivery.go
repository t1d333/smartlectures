package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/t1d333/smartlectures/internal/errors"
	"github.com/t1d333/smartlectures/internal/models"
	"github.com/t1d333/smartlectures/internal/notes"
	"github.com/t1d333/smartlectures/pkg/logger"
)

type Delivery struct {
	logger  logger.Logger
	service notes.Service
	mux     *gin.Engine
}

func NewDelivery(logger logger.Logger, mux *gin.Engine, service notes.Service) *Delivery {
	return &Delivery{
		logger:  logger,
		mux:     mux,
		service: service,
	}
}

func (d *Delivery) RegisterHandler() {
	d.mux.GET("/api/v1/notes/:noteId", d.GetNote)
	d.mux.GET("/api/v1/notes/overview", d.GetNotesOverview)
	d.mux.POST("/api/v1/notes", d.CreateNote)
	d.mux.DELETE("/api/v1/notes/:noteId", d.DeleteNote)
	d.mux.PUT("/api/v1/notes/:noteId", d.UpdateNote)
	d.mux.POST("/api/v1/notes/search", d.Search)
}

func (d *Delivery) GetNote(c *gin.Context) {
	noteIdStr := c.Param("noteId")

	if noteId, err := strconv.Atoi(noteIdStr); err != nil {
		_ = c.Error(errors.BadRequestError)
		return
	} else {
		note, err := d.service.GetNote(noteId, c.Request.Context())
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusOK, note)

	}
}

func (d *Delivery) Search(c *gin.Context) {
	req := models.SearchRequest{}
	if err := c.BindJSON(&req); err != nil {
		d.logger.Errorf("failed to decode search request: %w", err)
		_ = c.Error(errors.BadRequestError)
		return
	}

	result, err := d.service.SearchNote(req, c.Request.Context())
	if err != nil {
		d.logger.Errorf("failed to search note %s", err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.NoteSearchResult{
		Items: result,
	})
}

func (d *Delivery) CreateNote(c *gin.Context) {
	note := models.Note{}

	if err := c.BindJSON(&note); err != nil {
		d.logger.Errorf("failed to decode note: %w", err)
		_ = c.Error(errors.BadRequestError)
		return
	}

	noteId, err := d.service.CreateNote(note, c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"noteId": noteId})
}

func (d *Delivery) DeleteNote(c *gin.Context) {
	noteIdStr := c.Param("noteId")

	if noteId, err := strconv.Atoi(noteIdStr); err != nil {
		_ = c.Error(errors.BadRequestError)
		return
	} else if err := d.service.DeleteNote(noteId, c.Request.Context()); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (d *Delivery) GetNotesOverview(c *gin.Context) {
	// mock
	userId := 1

	overview, err := d.service.GetNotesOverview(userId, c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, overview)
}

func (d *Delivery) UpdateNote(c *gin.Context) {
	note := models.Note{}
	noteIdStr := c.Param("noteId")

	noteId := 0
	var err error

	if noteId, err = strconv.Atoi(noteIdStr); err != nil {
		_ = c.Error(errors.BadRequestError)
		return
	}

	if err := c.BindJSON(&note); err != nil {
		d.logger.Errorf("failed to decode note: %w", err)
		_ = c.Error(errors.BadRequestError)
		return
	}

	note.NoteId = noteId

	if err := d.service.UpdateNote(note, c.Request.Context()); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
