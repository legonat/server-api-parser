package handler

import (
	"awesomeProjectRucenter/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	r := gin.Default()

	authorized := r.Group("/")
	authorized.Use(gin.BasicAuth(gin.Accounts{"admin": "admin"}))
	{
		authorized.GET("/disks/", h.GetDisksWithLimit)
		authorized.GET("/vms/", h.GetVmsWithLimit)
	}
	return r
}
