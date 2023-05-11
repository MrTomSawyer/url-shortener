package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/MrTomSawyer/url-shortener/cmd/shortener/config"
)

type URLservice struct {
	repo   map[string]string
	config config.AppConfig
}

func (u *URLservice) ShortenURL(body string) string {
	hasher := md5.New()
	hasher.Write([]byte(body))
	hash := hex.EncodeToString(hasher.Sum(nil))[:8]
	shortURL := fmt.Sprintf("%s/%s", u.config.Server.DefaultAddr, hash)
	if _, ok := u.repo[hash]; !ok {
		u.repo[hash] = body
	}
	return shortURL
}

func (u *URLservice) GetOriginalURL(path string) (string, error) {
	if value, ok := u.repo[path]; ok {
		return value, nil
	} else {
		return "", fmt.Errorf("URL path '%s' not found", path)
	}
}
