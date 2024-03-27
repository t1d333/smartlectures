package http

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/t1d333/smartlectures/internal/errors"
	"github.com/t1d333/smartlectures/internal/recognizer"
	"github.com/t1d333/smartlectures/pkg/logger"
)

type Delivery struct {
	logger  logger.Logger
	service recognizer.Service
	mux     *gin.Engine
}

func NewDelivery(logger logger.Logger, mux *gin.Engine, service recognizer.Service) *Delivery {
	return &Delivery{
		logger:  logger,
		mux:     mux,
		service: service,
	}
}

func (d *Delivery) RegisterHandler() {
	d.mux.POST("/api/v1/recognizer/formula", d.RecognizeFormula)
	d.mux.POST("/api/v1/notes/recognizer/text", d.RecognizeText)
	d.mux.POST("/api/v1/notes/recognizer/mixed", d.RecognizeMixed)
}

func (d *Delivery) RecognizeFormula(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["image"]

	if len(files) != 1 {
		c.Error(errors.BadRequestError)
		return
	}

	img, _ := files[0].Open()
	defer img.Close()

	data, _ := io.ReadAll(img)

	formula, err := d.service.RecognizeFormula(data, c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"formula": formula})
}

func (d *Delivery) RecognizeText(c *gin.Context) {
}

func (d *Delivery) RecognizeMixed(c *gin.Context) {
}
