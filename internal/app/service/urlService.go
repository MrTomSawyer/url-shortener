package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/MrTomSawyer/url-shortener/internal/app/config"
	"github.com/MrTomSawyer/url-shortener/internal/app/repository"
)

type urlService struct {
	repo   repository.RepoHandler
	config config.AppConfig
}

func (u *urlService) ShortenURL(body string) (string, error) {
	hasher := md5.New()
	hasher.Write([]byte(body))
	hash := hex.EncodeToString(hasher.Sum(nil))[:8]

	val, err := u.repo.OriginalURL(hash)
	if err != nil {
		return "", fmt.Errorf("failed to check if such short url value presents: %v", err)
	}
	if val != "" {
		counter := 1
		for {
			newHash := hash + strconv.Itoa(counter)
			val, err := u.repo.OriginalURL(newHash)
			if err != nil {
				return "", fmt.Errorf("failed to check if such short url value presents: %v", err)
			}
			if val == "" {
				hash = newHash
				break
			}
			counter++
		}
	}

	err = u.repo.Create(hash, body)
	if err != nil {
		return "", fmt.Errorf("failed to save short URL: %v", err)
	}
	return fmt.Sprintf("%s/%s", u.config.Server.DefaultAddr, hash), nil
}

func (u *urlService) ExpandURL(path string) (string, error) {
	url, err := u.repo.OriginalURL(path)
	if err != nil {
		return "", fmt.Errorf("URL path '%s' not found", path)
	}
	return url, nil
}
