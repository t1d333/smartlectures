package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
	d.mux.POST("/api/v1/notes/:noteId", d.UpdateNote)
}

func (d *Delivery) GetNote(c *gin.Context) {
	noteIdStr := c.Param("noteId")

	if noteId, err := strconv.Atoi(noteIdStr); err != nil {
		c.Status(http.StatusBadRequest)
		return
	} else {
		note, err := d.service.GetNote(noteId, c.Request.Context())
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, note)

	}
}

func (d *Delivery) CreateNote(c *gin.Context) {
	note := models.Note{}

	if err := c.BindJSON(&note); err != nil {
		d.logger.Errorf("failed to decode note: %w", err)
		c.Status(http.StatusBadRequest)
		return
	}

	noteId, err := d.service.CreateNote(note, c.Request.Context())
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"noteId": noteId})
	c.Status(http.StatusNoContent)
}

func (d *Delivery) DeleteNote(c *gin.Context) {
	noteIdStr := c.Param("noteId")

	if noteId, err := strconv.Atoi(noteIdStr); err != nil {
		c.Status(http.StatusBadRequest)
		return
	} else if err := d.service.DeleteNote(noteId, c.Request.Context()); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

func (d *Delivery) GetNotesOverview(c *gin.Context) {
	// Is mock
	userId := 1

	overview, err := d.service.GetNotesOverview(userId, c.Request.Context())
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, overview)
}

func (d *Delivery) UpdateNote(c *gin.Context) {
	note := models.Note{}
	noteIdStr := c.Param("noteId")

	if noteId, err := strconv.Atoi(noteIdStr); err != nil {
		c.Status(http.StatusBadRequest)
		return
	} else {
		note.NoteId = noteId
	}

	if err := c.BindJSON(&note); err != nil {
		d.logger.Errorf("failed to decode note: %w", err)
		c.Status(http.StatusBadRequest)
		return
	}

	if err := d.service.UpdateNote(note, c.Request.Context()); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
