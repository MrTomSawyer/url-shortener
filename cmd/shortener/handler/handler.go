package handler

import (
	"github.com/MrTomSawyer/url-shortener/cmd/shortener/service"
	"github.com/gin-gonic/gin"
)

type logger interface {
	Infof(template string, args ...interface{})
}

type Handler struct {
	services *service.ServiceContainer
}

func NewHandler(services *service.ServiceContainer) *Handler {
	return &Handler{
		services: services,
	}
}

func (h Handler) InitRoutes(lg logger) *gin.Engine {
	router := gin.New()
	router.Use(h.logReqResInfo(lg))

	router.POST("/", h.ShortenURL)
	router.GET("/:id", h.ExpandURL)

	return router
}
