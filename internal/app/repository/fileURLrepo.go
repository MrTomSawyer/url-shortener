package repository

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"

	"github.com/MrTomSawyer/url-shortener/internal/app/logger"
	"github.com/MrTomSawyer/url-shortener/internal/app/models"
)

type FileURLrepo struct {
	storage     map[string]string
	path        string
	largestUUID int
}

func NewFileURLrepo(path string) (*FileURLrepo, error) {
	fileRepo := FileURLrepo{
		storage: map[string]string{},
		path:    path,
	}

	if err := fileRepo.Read(); err != nil {
		return nil, err
	}
	return &fileRepo, nil
}

func (s *FileURLrepo) Create(shortURL, originalURL, userID string) error {
	logger.Log.Infof("Writing to file... ShortURL: %s, OriginalURL: %s", shortURL, originalURL)

	uj := models.URLJson{
		UUID:        s.largestUUID + 1,
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}
	s.largestUUID += 1
	s.storage[shortURL] = originalURL

	file, err := os.OpenFile(s.path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("error opening file to write: %v", err)
	}
	defer file.Close()

	fileWriter := bufio.NewWriter(file)
	defer fileWriter.Flush()

	parsedURLJSON, err := json.Marshal(uj)
	if err != nil {
		return fmt.Errorf("error parsing data: %w", err)
	}

	_, err = fileWriter.Write(parsedURLJSON)
	if err != nil {
		return fmt.Errorf("error writing data to file: %v", err)
	}

	_, err = fileWriter.WriteString("\n")
	if err != nil {
		return fmt.Errorf("error writing a byte: %v", err)
	}

	return nil
}

func (s *FileURLrepo) Read() error {
	logger.Log.Infof("Starting to read file...")
	file, err := os.OpenFile(s.path, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return fmt.Errorf("error opening file to read: %v", err)
	}
	defer file.Close()

	var largestUUID int
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {

		var uj models.URLJson
		err := json.Unmarshal(fileScanner.Bytes(), &uj)
		if err != nil {
			return fmt.Errorf("error parsing line: %v", err)
		}

		url, err := url.Parse(uj.ShortURL)
		if err != nil {
			return err
		}
		shortURL := path.Base(url.Path)
		s.storage[shortURL] = uj.OriginalURL

		if uj.UUID > largestUUID {
			largestUUID = uj.UUID
		}

		s.largestUUID = largestUUID
	}
	logger.Log.Infof("File has been read")
	return nil
}

func (s FileURLrepo) OriginalURL(shortURL string) (string, error) {
	return s.storage[shortURL], nil
}

func (s FileURLrepo) BatchCreate(data []models.TempURLBatchRequest, userID string) ([]models.BatchURLResponce, error) {
	return []models.BatchURLResponce{}, nil
}

func (s FileURLrepo) GetAll(userid string) ([]models.URLJsonResponse, error) {
	var response []models.URLJsonResponse

	for key, value := range s.storage {
		response = append(response, models.URLJsonResponse{ShortURL: key, OriginalURL: value})
	}
	return response, nil
}
