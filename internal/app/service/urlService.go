package service

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"sync"

	"github.com/MrTomSawyer/url-shortener/internal/app/apperrors"
	"github.com/MrTomSawyer/url-shortener/internal/app/config"
	"github.com/MrTomSawyer/url-shortener/internal/app/logger"
	"github.com/MrTomSawyer/url-shortener/internal/app/models"
	"github.com/MrTomSawyer/url-shortener/internal/app/repository"
	"github.com/MrTomSawyer/url-shortener/internal/app/service/semaphore"
)

type urlService struct {
	Repo   repository.RepoHandler
	config config.AppConfig
}

func (u *urlService) ShortenURLHandler(body string, userID string) (string, error) {
	shortPath, err := u.ShortenURL(body)
	if err != nil {
		return "", err
	}

	err = u.Repo.Create(shortPath, body, userID)
	if err != nil {
		var urlConflictError *apperrors.URLConflict
		if errors.As(err, &urlConflictError) {
			return fmt.Sprintf("%s/%s", u.config.Server.DefaultAddr, urlConflictError.Value), err
		}
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
	if url == "" {
		return "", apperrors.ErrNotFound
	}
	return url, nil
}

func (u *urlService) HandleBatchInsert(data io.ReadCloser, userID string) ([]models.BatchURLResponce, error) {
	var parsedReq []models.BatchURLRequest

	decoder := json.NewDecoder(data)
	err := decoder.Decode(&parsedReq)
	if err != nil {
		logger.Log.Infof("Failed to decode json")
		return []models.BatchURLResponce{}, err
	}

	var tempURLRequests []models.TempURLBatchRequest

	for _, req := range parsedReq {
		shortURL, err := u.ShortenURL(req.OriginalURL)
		if err != nil {
			logger.Log.Infof("Failed to shorten URL")
			continue
		}
		tempURLRequests = append(tempURLRequests, models.TempURLBatchRequest{
			CorrelationID: req.CorrelationID,
			ShortURL:      shortURL,
			OriginalURL:   req.OriginalURL,
		})
	}

	switch {
	case u.config.DataBase.ConnectionStr != "":
		res, err := u.Repo.BatchCreate(tempURLRequests, userID)
		if err != nil {
			return []models.BatchURLResponce{}, err
		}
		return res, nil
	default:
		for _, req := range tempURLRequests {
			err := u.Repo.Create(req.ShortURL, req.OriginalURL, userID)
			if err != nil {
				return []models.BatchURLResponce{}, err
			}

		}
		return []models.BatchURLResponce{}, nil
	}
}

func (u *urlService) GetAll(userid string) ([]models.URLJsonResponse, error) {
	urls, err := u.Repo.GetAll(userid)
	if err != nil {
		return []models.URLJsonResponse{}, nil
	}

	return urls, nil
}

func (u *urlService) DeleteAll(urls []string, userid string) {
	var wg sync.WaitGroup
	semaphore := semaphore.NewSemaphore(2)

	resultCh := make(chan error, len(urls))

	for _, url := range urls {
		wg.Add(1)

		go func(url string, userid string, resultCh chan error) {
			semaphore.Acquire()
			defer wg.Done()
			defer semaphore.Release()

			err := u.Repo.DeleteAll(url, userid)
			if err != nil {
				resultCh <- err
			}

		}(url, userid, resultCh)
	}

	for err := range resultCh {
		fmt.Printf("error deleting url: %v", err)
	}

	wg.Wait()
}
