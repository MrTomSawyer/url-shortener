package repository

import "github.com/MrTomSawyer/url-shortener/internal/app/models"

type InMemoryURLRepo struct {
	storage map[string]string
}

func NewInMemoryURLRepo() *InMemoryURLRepo {
	return &InMemoryURLRepo{
		storage: map[string]string{},
	}
}

func (i *InMemoryURLRepo) Create(shortURL, originalURL string) error {
	i.storage[shortURL] = originalURL
	return nil
}

func (i *InMemoryURLRepo) OriginalURL(shortURL string) (string, error) {
	return i.storage[shortURL], nil
}

func (i *InMemoryURLRepo) BatchCreate(data []models.TempURLBatchRequest) ([]models.BatchURLResponce, error) {
	return []models.BatchURLResponce{}, nil
}
