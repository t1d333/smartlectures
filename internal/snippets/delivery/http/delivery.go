package http

import (
	"net/http"
	// "strconv"

	"github.com/gin-gonic/gin"
	// "github.com/t1d333/smartlectures/internal/errors"
	// "github.com/t1d333/smartlectures/internal/models"
	"github.com/t1d333/smartlectures/internal/snippets"
	"github.com/t1d333/smartlectures/pkg/logger"
)

type Delivery struct {
	logger  logger.Logger
	service snippets.Service
	mux     *gin.Engine
}

func NewDelivery(logger logger.Logger, mux *gin.Engine, service snippets.Service) *Delivery {
	return &Delivery{
		logger:  logger,
		mux:     mux,
		service: service,
	}
}

func (d *Delivery) RegisterHandler() {
	d.mux.GET("/api/v1/snippets", d.GetSnippets)
	d.mux.POST("/api/v1/snippets", d.CreateSnippet)
	d.mux.DELETE("/api/v1/snippets/:snippetId", d.DeleteSnippet)
	d.mux.PUT("/api/v1/snippets/:snippetId", d.UpdateSnippet)
	d.mux.POST("/api/v1/snippets/search", d.UpdateSnippet)
}

func (d *Delivery) GetSnippets(c *gin.Context) {
	userId := 1

	snippets, err := d.service.GetSnippets(userId, c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"snippets": snippets})
}

func (d *Delivery) Search(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func (d *Delivery) CreateSnippet(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func (d *Delivery) DeleteSnippet(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func (d *Delivery) UpdateSnippet(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
