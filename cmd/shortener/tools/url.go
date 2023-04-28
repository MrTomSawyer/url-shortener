package tools

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func Shorten(url string) (string, string) {
	hasher := md5.New()
	hasher.Write([]byte(url))
	hash := hex.EncodeToString(hasher.Sum(nil))[:8]
	return fmt.Sprintf("http://localhost:8080/%s", hash), hash
}
