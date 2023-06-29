package handler

import (
	"errors"
	"net/http"

	"github.com/MrTomSawyer/url-shortener/internal/app/apperrors"
	"github.com/gin-gonic/gin"
)

func (h *Handler) ExpandURL(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		http.Error(c.Writer, "Error reading id param", http.StatusInternalServerError)
		return
	}

	value, err := h.services.URL.ExpandURL(id)
	if err != nil {
		if errors.Is(err, apperrors.ErrURLDeleted) {
			http.Error(c.Writer, "URL has been deleted", http.StatusGone)
			return
		}
		http.Error(c.Writer, "No original URL found", http.StatusNotFound)
		return
	}

	c.Status(http.StatusTemporaryRedirect)
	c.Header("Location", value)
}
