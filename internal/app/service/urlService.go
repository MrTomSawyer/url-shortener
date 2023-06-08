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
	Repo   repository.RepoHandler
	config config.AppConfig
}

func (u *urlService) ShortenURLHandler(body string) (string, error) {
	shortPath, err := u.ShortenURL(body)
	if err != nil {
		return "", err
	}

	err = u.Repo.Create(shortPath, body)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", u.config.Server.DefaultAddr, shortPath), nil
}

func (u *urlService) ShortenURL(body string) (string, error) {
	hasher := md5.New()
	hasher.Write([]byte(body))
	shortPath := hex.EncodeToString(hasher.Sum(nil))[:8]

	val, err := u.Repo.OriginalURL(shortPath)
	if err != nil {
		return "", fmt.Errorf("failed to check if such short url value presents: %v", err)
	}
	if val != "" {
		counter := 1
		for {
			newShortPath := shortPath + strconv.Itoa(counter)
			val, err := u.Repo.OriginalURL(newShortPath)
			if err != nil {
				return "", fmt.Errorf("failed to check if such short url value presents: %v", err)
			}
			if val == "" {
				shortPath = newShortPath
				break
			}
			counter++
		}
	}
	return shortPath, nil
}

func (u *urlService) ExpandURL(path string) (string, error) {
	url, err := u.Repo.OriginalURL(path)
	if err != nil {
		return "", fmt.Errorf("URL path '%s' not found", path)
	}
	return url, nil
}
