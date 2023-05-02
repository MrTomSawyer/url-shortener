package tools

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"

	f "github.com/MrTomSawyer/url-shortener/cmd/shortener/config"
)

var addr string

func Shorten(url string) (string, string) {
	hasher := md5.New()
	hasher.Write([]byte(url))
	hash := hex.EncodeToString(hasher.Sum(nil))[:8]
	if f.DefaultAddr == "" {
		addr = "http://localhost:8080"
	}
	return fmt.Sprintf("%s/%s", addr, hash), hash
}
