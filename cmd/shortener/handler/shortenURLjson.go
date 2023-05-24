package handler

import (
	"encoding/json"
	"net/http"

	"github.com/MrTomSawyer/url-shortener/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) shortenURLjson(c *gin.Context) {
	var req models.ShortenRequest

	dec := json.NewDecoder(c.Request.Body)
	if err := dec.Decode(&req); err != nil {
		http.Error(c.Writer, "Error decoding JSON body", http.StatusInternalServerError)
		return
	}

	shortenURL := h.services.URL.ShortenURL(req.Url)

	res := models.ShortenResponce{
		Result: shortenURL,
	}

	c.JSON(http.StatusOK, res)
}
