package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/t1d333/smartlectures/internal/dirs"
	"github.com/t1d333/smartlectures/internal/errors"
	"github.com/t1d333/smartlectures/internal/models"
	"github.com/t1d333/smartlectures/pkg/logger"
)

type Delivery struct {
	logger  logger.Logger
	service dirs.Service
	mux     *gin.Engine
}

func NewDelivery(logger logger.Logger, mux *gin.Engine, service dirs.Service) *Delivery {
	return &Delivery{
		logger:  logger,
		mux:     mux,
		service: service,
	}
}

func (d *Delivery) RegisterHandler() {
	d.mux.GET("/api/v1/dirs/:dirId", d.GetDir)
	d.mux.GET("/api/v1/dirs/overview", d.GetDirsOverview)
	d.mux.POST("/api/v1/dirs", d.CreateDir)
	d.mux.DELETE("/api/v1/dirs/:dirId", d.DeleteDir)
	d.mux.PUT("/api/v1/dirs/:dirId", d.UpdateDir)
}

func (d *Delivery) GetDir(c *gin.Context) {
	dirIdStr := c.Param("dirId")

	if dirId, err := strconv.Atoi(dirIdStr); err != nil {
		_ = c.Error(errors.BadRequestError)
		return
	} else {
		dir, err := d.service.GetDir(dirId, c.Request.Context())
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusOK, dir)

	}
}

func (d *Delivery) CreateDir(c *gin.Context) {
	dir := models.Dir{}

	if err := c.BindJSON(&dir); err != nil {
		d.logger.Errorf("failed to decode dir: %w", err)
		_ = c.Error(errors.BadRequestError)
		return
	}

	dirId, err := d.service.CreateDir(dir, c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"dirId": dirId})
}

func (d *Delivery) DeleteDir(c *gin.Context) {
	dirIdStr := c.Param("dirId")

	if dirId, err := strconv.Atoi(dirIdStr); err != nil {
		_ = c.Error(errors.BadRequestError)
		return
	} else if err := d.service.DeleteDir(dirId, c.Request.Context()); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (d *Delivery) GetDirsOverview(c *gin.Context) {
	userId := 1

	overview, err := d.service.GetDirsOverview(userId, c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, overview)
}

func (d *Delivery) UpdateDir(c *gin.Context) {
	dir := models.Dir{}
	dirIdStr := c.Param("dirId")
	dirId := 0
	var err error

	if dirId, err = strconv.Atoi(dirIdStr); err != nil {
		_ = c.Error(errors.BadRequestError)
		return
	}

	if err := c.BindJSON(&dir); err != nil {
		d.logger.Errorf("failed to decode dir: %w", err)
		_ = c.Error(errors.BadRequestError)
		return
	}

	dir.DirId = dirId

	if err := d.service.UpdateDir(dir, c.Request.Context()); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
