package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"

	"github.com/MrTomSawyer/url-shortener/internal/app/models"
)

type Storage struct {
	path        string
	largestUUID int
}

func (s *Storage) Write(uj *models.URLJson) error {
	file, err := os.OpenFile(s.path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("error opening file to write: %v", err)
	}
	defer file.Close()

	fileWriter := bufio.NewWriter(file)
	defer fileWriter.Flush()

	parsedURLJSON, err := json.Marshal(uj)
	if err != nil {
		return fmt.Errorf("error parsing data: %v", err)
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

func (s *Storage) Read(repo *map[string]string) error {
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
			panic(err)
		}
		shortURL := path.Base(url.Path)
		(*repo)[shortURL] = uj.OriginalURL

		if uj.UUID > largestUUID {
			largestUUID = uj.UUID
		}

		s.largestUUID = largestUUID
	}

	return nil
}

func NewStorage(path string) (*Storage, error) {
	_, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	return &Storage{
		path: path,
	}, nil
}