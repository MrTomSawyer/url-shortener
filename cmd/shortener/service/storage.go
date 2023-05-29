package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Storage struct {
	file   *os.File
	path   string
	reader bufio.Reader
	writer bufio.Writer
}

func (s *Storage) Write(p []byte) (int, error) {
	data, err := json.Marshal(p)
	if err != nil {
		return 0, fmt.Errorf("error marshalling data")
	}

	if _, err := s.writer.Write(p); err != nil {
		return 0, fmt.Errorf("error writing data")
	}

	if err := s.writer.WriteByte('\n'); err != nil {
		return 0, fmt.Errorf("error writing data")
	}

	return len(data), s.writer.Flush()
}

func (s *Storage) Read(p []byte) (int, error) {
	n, err := io.ReadFull(s.file, p)
	if err != nil {
		return n, fmt.Errorf("error reading file: %v", err)
	}
	return n, err
}

func NewStorage(path string) (*Storage, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			currentDir, err := os.Getwd()
			if err != nil {
				return nil, err
			}
			filePath := filepath.Join(currentDir, path)
			file, err := os.Create(filePath)
			if err != nil {
				return nil, fmt.Errorf("error creating file: %v", err)
			}
			defer file.Close()

			return &Storage{
				file:   file,
				path:   path,
				reader: *bufio.NewReader(file),
				writer: *bufio.NewWriter(file),
			}, nil
		}
		return nil, err
	}
	return nil, fmt.Errorf("this storage already exists")
}
