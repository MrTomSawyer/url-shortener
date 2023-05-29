package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	m "github.com/MrTomSawyer/url-shortener/internal/models"
)

type Storage struct {
	file   *os.File
	reader bufio.Reader
	writer bufio.Writer
}

func (s *Storage) Write(uj *m.URLJson) error {
	data, err := json.Marshal(uj)
	if err != nil {
		return err
	}

	if _, err := s.writer.Write(data); err != nil {
		return err
	}

	if err := s.writer.WriteByte('\n'); err != nil {
		return err
	}

	return s.writer.Flush()
}

func (s *Storage) Read() (*m.URLJson, error) {
	data, err := s.reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	uj := m.URLJson{}
	err = json.Unmarshal(data, &uj)
	if err != nil {
		return nil, err
	}
	return &uj, nil
}

func (s Storage) LastUUID() (int, error) {
	largestUUID := 0

	_, err := s.file.Seek(0, 0)
	if err != nil {
		return 0, err
	}

	for {
		line, err := s.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return 0, err
		}

		uj := m.URLJson{}
		err = json.Unmarshal([]byte(line), &uj)
		if err != nil {
			return 0, err
		}

		if uj.UUID > largestUUID {
			largestUUID = uj.UUID
		}
	}
	return largestUUID, nil
}

func NewStorage(path string) (*Storage, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	filePath := filepath.Join(currentDir, path)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, fmt.Errorf("error creating file: %v", err)
	}

	return &Storage{
		file:   file,
		reader: *bufio.NewReader(file),
		writer: *bufio.NewWriter(file),
	}, nil
}
