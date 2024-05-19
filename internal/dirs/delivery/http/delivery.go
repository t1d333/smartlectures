package http

import (
	"context"
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
	userId := c.Keys["userId"].(int)

	if dirId, err := strconv.Atoi(dirIdStr); err != nil {
		_ = c.Error(errors.BadRequestError)
		return
	} else {
		dir, err := d.service.GetDir(context.WithValue(c.Request.Context(), "userId", userId), dirId)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusOK, dir)

	}
}

func (d *Delivery) CreateDir(c *gin.Context) {
	dir := models.Dir{}
	userId := c.Keys["userId"].(int)

	if err := c.BindJSON(&dir); err != nil {
		d.logger.Errorf("failed to decode dir: %w", err)
		_ = c.Error(errors.BadRequestError)
		return
	}

	dirId, err := d.service.CreateDir(context.WithValue(c.Request.Context(), "userId", userId), dir)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"dirId": dirId})
}

func (d *Delivery) DeleteDir(c *gin.Context) {
	dirIdStr := c.Param("dirId")
	userId := c.Keys["userId"].(int)

	if dirId, err := strconv.Atoi(dirIdStr); err != nil {
		_ = c.Error(errors.BadRequestError)
		return
	} else if err := d.service.DeleteDir(context.WithValue(c.Request.Context(), "userId", userId), dirId); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (d *Delivery) GetDirsOverview(c *gin.Context) {
	userId := c.Keys["userId"].(int)
	
	overview, err := d.service.GetDirsOverview(c.Request.Context(), userId)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, overview)
}

func (d *Delivery) UpdateDir(c *gin.Context) {
	userId := c.Keys["userId"].(int)
	
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

	if err := d.service.UpdateDir(context.WithValue(c.Request.Context(), "userId", userId), dir); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
