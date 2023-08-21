// Package handler provides HTTP request handlers for managing URL-related operations.
package handler

import (
	"github.com/MrTomSawyer/url-shortener/internal/app/config"
	"github.com/MrTomSawyer/url-shortener/internal/app/middlewares"
	"github.com/MrTomSawyer/url-shortener/internal/app/service"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

// Handler is responsible for defining and initializing HTTP request handlers.
type Handler struct {
	services *service.ServiceContainer
	Cfg      config.AppConfig
}

// NewHandler creates a new instance of Handler.
func NewHandler(services *service.ServiceContainer, cfg config.AppConfig) *Handler {
	return &Handler{
		services: services,
		Cfg:      cfg,
	}
}

// InitRoutes initializes the routes for the application and returns a gin.Engine instance.
func (h Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	pprof.Register(router)

	router.Use(
		middlewares.LogReqResInfo(),
		middlewares.DataCompressor(),
		middlewares.CookieHandler(h.Cfg.Server.SecretKey),
	)

	router.POST("/", h.ShortenURL)
	router.GET("/:id", h.ExpandURL)
	router.GET("/ping", h.pingDB)

	api := router.Group("/api")
	{
		api.POST("/shorten", h.ShortenURLjson)
		api.POST("/shorten/batch", h.batchURLinsert)
		api.GET("/user/urls", h.GetAll)
		api.DELETE("/user/urls", h.deleteAll)
	}

	return router
}
