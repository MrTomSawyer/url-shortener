package main

import (
	f "github.com/MrTomSawyer/url-shortener/cmd/shortener/config"
	h "github.com/MrTomSawyer/url-shortener/cmd/shortener/httphandlers"
	"github.com/gin-gonic/gin"
)

func main() {
	f.ParseFlags()
	app := gin.Default()

	app.POST("/", func(c *gin.Context) {
		h.ShortenURL(c.Writer, c.Request)
	})
	app.GET("/:id", func(c *gin.Context) {
		h.GetOriginalURL(c.Writer, c.Request)
	})

	err := app.Run(f.ServerAddr)
	if err != nil {
		panic(err)
	}
}
