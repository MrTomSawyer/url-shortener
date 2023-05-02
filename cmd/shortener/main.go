package main

import (
	h "github.com/MrTomSawyer/url-shortener/cmd/shortener/httphandlers"
	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()

	app.POST("/", func(c *gin.Context) {
		h.ShortenURL(c.Writer, c.Request)
	})
	app.GET("/:id", func(c *gin.Context) {
		h.GetOriginalURL(c.Writer, c.Request)
	})

	err := app.Run(`:8080`)
	if err != nil {
		panic(err)
	}
}
