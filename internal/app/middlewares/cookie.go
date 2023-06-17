package middlewares

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CookieHandler(secret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		h := hmac.New(sha256.New, []byte(secret))

		cookie, err := ctx.Cookie("user_id")
		if err != nil {
			c := createUUIDCookie(h)
			ctx.SetCookie("user_id", c, 0, "/", "", false, true)
			ctx.Next()
		}

		if cookie == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		values := strings.Split(cookie, "|")
		if len(values) != 2 {
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
