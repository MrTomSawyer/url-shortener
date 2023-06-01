package handler

import (
	"github.com/MrTomSawyer/url-shortener/internal/app/middlewares"
	"github.com/MrTomSawyer/url-shortener/internal/app/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.ServiceContainer
}

func NewHandler(services *service.ServiceContainer) *Handler {
	return &Handler{
		services: services,
	}
}

func (h Handler) InitRoutes(lg middlewares.Logger) *gin.Engine {
	router := gin.New()
	router.Use(middlewares.LogReqResInfo(lg), middlewares.DataCompressor())

	router.POST("/", h.ShortenURL)
	router.GET("/:id", h.ExpandURL)

	api := router.Group("/api")
	{
		api.POST("/shorten", h.ShortenURLjson)
	}

	return router
}
