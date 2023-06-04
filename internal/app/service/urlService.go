package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/MrTomSawyer/url-shortener/internal/app/config"
	"github.com/MrTomSawyer/url-shortener/internal/app/models"
	"github.com/MrTomSawyer/url-shortener/internal/app/repository"
)

type urlService struct {
	repo    map[string]string
	db      *repository.RepositoryContainer
	config  config.AppConfig
	storage *Storage
}

func (u *urlService) ShortenURL(body string) (string, error) {
	hasher := md5.New()
	hasher.Write([]byte(body))
	hash := hex.EncodeToString(hasher.Sum(nil))[:8]

	var shortURL string
	if _, ok := u.repo[hash]; !ok {
		u.repo[hash] = body
		shortURL = fmt.Sprintf("%s/%s", u.config.Server.DefaultAddr, hash)
	} else {
		counter := 1
		for {
			newHash := hash + strconv.Itoa(counter)
			if _, ok := u.repo[newHash]; !ok {
				u.repo[newHash] = body
				shortURL = fmt.Sprintf("%s/%s", u.config.Server.DefaultAddr, newHash)
				break
			}
			counter++
		}
	}

	uj := models.URLJson{
		UUID:        u.storage.largestUUID + 1,
		ShortURL:    shortURL,
		OriginalURL: body,
	}
	err := u.storage.Write(&uj)
	if err != nil {
		fmt.Println("Failed to write data to file", err)
	}
	u.storage.largestUUID++

	return shortURL, nil
}

func (u *urlService) ExpandURL(path string) (string, error) {
	if value, ok := u.repo[path]; ok {
		return value, nil
	} else {
		return "", fmt.Errorf("URL path '%s' not found", path)
	}
}
