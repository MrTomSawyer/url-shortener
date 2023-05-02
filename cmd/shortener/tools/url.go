package tools

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"

	f "github.com/MrTomSawyer/url-shortener/cmd/shortener/config"
)

func Shorten(url string) (string, string) {
	hasher := md5.New()
	hasher.Write([]byte(url))
	hash := hex.EncodeToString(hasher.Sum(nil))[:8]
	return fmt.Sprintf("%s/%s", f.DefaultAddr, hash), hash
}
