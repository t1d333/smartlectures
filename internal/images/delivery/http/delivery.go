package http

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/t1d333/smartlectures/internal/errors"
	"github.com/t1d333/smartlectures/internal/images"
	// "github.com/t1d333/smartlectures/internal/models"
	"github.com/t1d333/smartlectures/pkg/logger"
)

type Delivery struct {
	logger  logger.Logger
	service images.Service
	mux     *gin.Engine
}

func NewDelivery(logger logger.Logger, mux *gin.Engine, service images.Service) *Delivery {
	return &Delivery{
		logger:  logger,
		mux:     mux,
		service: service,
	}
}

func (d *Delivery) RegisterHandler() {
	d.mux.POST("/api/v1/images/upload", d.UploadImage)
}

func (d *Delivery) UploadImage(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["image"]

	if len(files) < 1 {
		_ = c.Error(errors.BadRequestError)
		return
	}

	img, _ := files[0].Open()
	defer img.Close()

	data, _ := io.ReadAll(img)
	reader := bytes.NewBuffer(data)

	src, err := d.service.UploadImage(reader, c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"src": src})
}
