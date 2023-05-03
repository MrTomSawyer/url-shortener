package handler

import (
	"github.com/MrTomSawyer/url-shortener/cmd/shortener/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	root := router.Group("/")
	{
		root.POST("/", h.ShortenURL)
		root.GET("/:id", h.GetOriginalURL)
	}

	return router
}
