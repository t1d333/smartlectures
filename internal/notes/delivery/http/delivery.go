package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	c.String(http.StatusOK, c.Request.URL.Path)
}

func (d *Delivery) CreateNote(c *gin.Context) {
	c.String(http.StatusOK, c.Request.URL.Path)
}

func (d *Delivery) GetNotesOverview(c *gin.Context) {
	c.String(http.StatusOK, c.Request.URL.Path)
}

func (d *Delivery) DeleteNote(c *gin.Context) {
	c.String(http.StatusOK, c.Request.URL.Path)
}

func (d *Delivery) UpdateNote(c *gin.Context) {
	c.String(http.StatusOK, c.Request.URL.Path)
}
