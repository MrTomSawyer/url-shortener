// Package middlewares provides middleware functions for handling user authentication and cookies.
package middlewares

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"net/http"
	"strings"
	"time"

	"github.com/MrTomSawyer/url-shortener/internal/app/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CookieHandler is a middleware that manages user authentication using cookies.
func CookieHandler(secret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		h := hmac.New(sha256.New, []byte(secret))
		maxAge := time.Now().Add(24 * time.Hour)

		cookie, err := ctx.Cookie("user_id")
		if err != nil {
			logger.Log.Info("No cookie found. Creating a new one...")
			c, userID := createUUIDCookie(h)
			ctx.SetCookie("user_id", c, int(maxAge.Unix()), "/", "", false, true)
			ctx.Set("user_id", userID)
			ctx.Next()
			return
		}

		values := strings.Split(cookie, "|")
		logger.Log.Infof("Cookie user ID and signature: %v", values)
		if len(values) != 2 {
			logger.Log.Info("Invalid cookie format")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userID := values[0]
		signature := values[1]

		h.Write([]byte(userID))
		expectedSignature := hex.EncodeToString(h.Sum(nil))

		if signature != expectedSignature {
			logger.Log.Info("Invalid cookie signature. Creating a new one...")
			c, userID := createUUIDCookie(h)
			ctx.SetCookie("user_id", c, int(maxAge.Unix()), "/", "", false, true)
			ctx.Set("user_id", userID)
			ctx.Next()
			return
		}

		ctx.Set("user_id", userID)
		ctx.Next()
	}
}

// createUUIDCookie generates a new UUID, calculates its HMAC signature, and returns both the concatenated cookie value and the UUID string.
func createUUIDCookie(h hash.Hash) (string, string) {
	uuid := uuid.New()
	h.Write([]byte(uuid.String()))
	signature := hex.EncodeToString(h.Sum(nil))
	return uuid.String() + "|" + signature, uuid.String()
}
