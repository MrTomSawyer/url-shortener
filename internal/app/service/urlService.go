// Package service provides functionalities for managing URL-related operations and business logic.
package service

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/MrTomSawyer/url-shortener/internal/app/apperrors"
	"github.com/MrTomSawyer/url-shortener/internal/app/config"
	"github.com/MrTomSawyer/url-shortener/internal/app/logger"
	"github.com/MrTomSawyer/url-shortener/internal/app/models"
	"github.com/MrTomSawyer/url-shortener/internal/app/repository"
)

// urlService implements the URL-related service methods.
type urlService struct {
	Repo   repository.RepoHandler
	config config.AppConfig
}

// ShortenURLHandler receives a URL, shortens it, and handles its insertion.
// Returns the shortened URL or an error if the operation fails.
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

// ShortenURL generates a shortened URL path based on an MD5 hash of the provided URL.
// Returns the shortened URL path or an error if the operation fails.
func (u *urlService) ShortenURL(body string) (string, error) {
	hasher := md5.New()
	hasher.Write([]byte(body))
	shortPath := hex.EncodeToString(hasher.Sum(nil))[:8]

	val, err := u.Repo.OriginalURL(shortPath)
	if err != nil {
		return "", fmt.Errorf("failed to check if such short URL value exists: %v", err)
	}
	if val != "" {
		counter := 1
		for {
			newShortPath := shortPath + strconv.Itoa(counter)
			val, err := u.Repo.OriginalURL(newShortPath)
			if err != nil {
				return "", fmt.Errorf("failed to check if such short URL value exists: %v", err)
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

// ExpandURL converts a shortened URL path back to its original URL.
// Returns the original URL or an error if the path is not found.
func (u *urlService) ExpandURL(path string) (string, error) {
	url, err := u.Repo.OriginalURL(path)
	if err != nil {
		return "", fmt.Errorf("URL path '%s' not found: %w", path, err)
	}
	if url == "" {
		return "", apperrors.ErrNotFound
	}
	return url, nil
}

// HandleBatchInsert processes a batch of URL insert requests.
// Returns a list of BatchURLResponce objects indicating the results of the batch insert.
func (u *urlService) HandleBatchInsert(data io.ReadCloser, userID string) ([]models.BatchURLResponse, error) {
	var parsedReq []models.BatchURLRequest

	decoder := json.NewDecoder(data)
	err := decoder.Decode(&parsedReq)
	if err != nil {
		logger.Log.Infof("Failed to decode JSON")
		return []models.BatchURLResponse{}, err
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
			return []models.BatchURLResponse{}, err
		}
		return res, nil
	default:
		for _, req := range tempURLRequests {
			err := u.Repo.Create(req.ShortURL, req.OriginalURL, userID)
			if err != nil {
				return []models.BatchURLResponse{}, err
			}

		}
		return []models.BatchURLResponse{}, nil
	}
}

// GetAll retrieves all URLs associated with a given user ID.
// Returns a list of URLJsonResponse objects or an error if the retrieval fails.
func (u *urlService) GetAll(userID string) ([]models.URLJsonResponse, error) {
	urls, err := u.Repo.GetAll(userID)
	if err != nil {
		return []models.URLJsonResponse{}, err
	}

	return urls, nil
}

// DeleteAll deletes multiple URLs associated with a given user ID.
// Returns an error if the deletion operation fails.
func (u *urlService) DeleteAll(urls []string, userID string) error {
	err := u.Repo.DeleteAll(urls, userID)
	if err != nil {
		fmt.Printf("failed to asynchronously delete URLs: %v", err)
		return err
	}
	return nil
}
