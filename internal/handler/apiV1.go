package handler

import (
	"awesomeProjectRucenter/pkg/erx"
	"awesomeProjectRucenter/pkg/tools"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

var log *logrus.Logger

func init() {
	log = tools.GetLogrusInstance("")
}

func (h *Handler) GetDisksWithLimit(c *gin.Context) {
	limit := c.DefaultQuery("limit", "25")
	offset := c.DefaultQuery("offset", "0")

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse query parameter"})
		log.Error(erx.New(err))
		return
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse query parameter"})
		log.Error(erx.New(err))
		return
	}

	res, err := h.service.GetDisksWithLimit(limitInt, offsetInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get data"})
		log.Error(erx.New(err))
		return
	}

	c.JSON(http.StatusOK, res)

}

func (h *Handler) GetVmsWithLimit(c *gin.Context) {
	limit := c.DefaultQuery("limit", "25")
	offset := c.DefaultQuery("offset", "0")

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse query parameter"})
		log.Error(erx.New(err))
		return
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse query parameter"})
		log.Error(erx.New(err))
		return
	}

	res, err := h.service.GetVmsWithLimit(limitInt, offsetInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get data"})
		log.Error(erx.New(err))
		return
	}

	c.JSON(http.StatusOK, res)

}
