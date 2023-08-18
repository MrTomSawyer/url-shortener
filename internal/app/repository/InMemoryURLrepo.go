// Package repository provides implementations for data storage and retrieval.
package repository

import "github.com/MrTomSawyer/url-shortener/internal/app/models"

// InMemoryURLRepo is an in-memory repository for storing URL data.
type InMemoryURLRepo struct {
	storage map[string]string
}

// NewInMemoryURLRepo creates a new instance of InMemoryURLRepo.
func NewInMemoryURLRepo() *InMemoryURLRepo {
	return &InMemoryURLRepo{
		storage: map[string]string{},
	}
}

// Create adds a new URL entry to the in-memory repository.
func (i *InMemoryURLRepo) Create(shortURL, originalURL, userID string) error {
	i.storage[shortURL] = originalURL
	return nil
}

// OriginalURL retrieves the original URL corresponding to a short URL from the in-memory repository.
func (i *InMemoryURLRepo) OriginalURL(shortURL string) (string, error) {
	return i.storage[shortURL], nil
}

// BatchCreate is not implemented for InMemoryURLRepo.
func (i *InMemoryURLRepo) BatchCreate(data []models.TempURLBatchRequest, userID string) ([]models.BatchURLResponse, error) {
	return []models.BatchURLResponse{}, nil
}

// GetAll retrieves all URLs for a given user from the in-memory repository.
func (i *InMemoryURLRepo) GetAll(userID string) ([]models.URLJsonResponse, error) {
	var response []models.URLJsonResponse

	for shortURL, originalURL := range i.storage {
		response = append(response, models.URLJsonResponse{ShortURL: shortURL, OriginalURL: originalURL})
	}
	return response, nil
}

// DeleteAll is not implemented for InMemoryURLRepo.
func (i *InMemoryURLRepo) DeleteAll(shortURLs []string, userID string) error {
	return nil
}
