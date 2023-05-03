package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/MrTomSawyer/url-shortener/cmd/shortener/config"
)

type URLservice struct {
	repo map[string]string
}

func NewURLservice(m map[string]string) *URLservice {
	return &URLservice{
		repo: make(map[string]string),
	}
}

func (u *URLservice) ShortenURL(body string) string {
	hasher := md5.New()
	hasher.Write([]byte(body))
	hash := hex.EncodeToString(hasher.Sum(nil))[:8]
	shortURL := fmt.Sprintf("%s/%s", config.DefaultAddr, hash)
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
