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

func CookieHandler(secret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		h := hmac.New(sha256.New, []byte(secret))
		logger.Log.Info("All cookies:", ctx.Request.Cookies())
		cookie, err := ctx.Cookie("user_id")
		logger.Log.Info("Cookie: %s", cookie)
		maxAge := time.Now().Add(24 * time.Hour)
		if err != nil {
			logger.Log.Info("Cookie err: %s", err)
			logger.Log.Info("Creating new cookie")
			c, userID := createUUIDCookie(h)

			ctx.SetCookie("user_id", c, int(maxAge.Unix()), "/", "", false, true)
			ctx.Set("user_id", userID)
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
			c, _ := createUUIDCookie(h)
			ctx.SetCookie("user_id", c, int(maxAge.Unix()), "/", "", false, true)
		}
		ctx.Set("user_id", userID)
		ctx.Next()
	}
}

func createUUIDCookie(h hash.Hash) (string, string) {
	uuid := uuid.New()
	h.Write([]byte(uuid.String()))
	signature := hex.EncodeToString(h.Sum(nil))
	return uuid.String() + "|" + signature, uuid.String()
}
