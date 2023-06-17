package middlewares

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"net/http"
	"strings"

	"github.com/MrTomSawyer/url-shortener/internal/app/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CookieHandler(secret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		h := hmac.New(sha256.New, []byte(secret))

		cookie, err := ctx.Cookie("user_id")
		logger.Log.Info("Cookie: ", cookie, err)
		if err != nil {
			logger.Log.Info("Creating new cookie")
			c := createUUIDCookie(h)
			ctx.SetCookie("user_id", c, 0, "/", "", false, true)
			return
		}

		if cookie == "" {
			logger.Log.Info("Cookie is empry string")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		values := strings.Split(cookie, "|")
		if len(values) != 2 {
			logger.Log.Info("Cookie includes less than 2 parts")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userID := values[0]
		signature := values[1]

		h.Write([]byte(userID))
		expectedSignature := hex.EncodeToString(h.Sum(nil))

		if signature != expectedSignature {
			c := createUUIDCookie(h)
			ctx.SetCookie("user_id", c, 0, "/", "", false, true)
		}

		ctx.Next()
	}
}

func createUUIDCookie(h hash.Hash) string {
	uuid := uuid.New()
	h.Write([]byte(uuid.String()))
	signature := hex.EncodeToString(h.Sum(nil))
	return uuid.String() + "|" + signature
}
