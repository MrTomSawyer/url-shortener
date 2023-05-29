package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/MrTomSawyer/url-shortener/cmd/shortener/config"
	m "github.com/MrTomSawyer/url-shortener/internal/models"
)

type urlService struct {
	repo     map[string]string
	config   config.AppConfig
	storage  *Storage
	lastUUID int
}

func (u *urlService) ShortenURL(body string) (string, error) {
	hasher := md5.New()
	hasher.Write([]byte(body))
	hash := hex.EncodeToString(hasher.Sum(nil))[:8]

	var shortURL string
	if _, ok := u.repo[hash]; !ok {
		u.repo[hash] = body
		shortURL = fmt.Sprintf("%s/%s", u.config.Server.DefaultAddr, hash)
	} else {
		counter := 1
		for {
			newHash := hash + strconv.Itoa(counter)
			if _, ok := u.repo[newHash]; !ok {
				u.repo[newHash] = body
				shortURL = fmt.Sprintf("%s/%s", u.config.Server.DefaultAddr, newHash)
				break
			}
			counter++
		}
	}

	uj := m.URLJson{
		UUID:        u.lastUUID + 1,
		ShortURL:    shortURL,
		OriginalURL: body,
	}
	err := u.storage.Write(&uj)
	if err != nil {
		fmt.Println("Failed to write data to file", err)
	}
	u.lastUUID++

	return shortURL, nil
}

func (u *urlService) ExpandURL(path string) (string, error) {
	if value, ok := u.repo[path]; ok {
		return value, nil
	} else {
		return "", fmt.Errorf("URL path '%s' not found", path)
	}
}

func (u *urlService) initializeLastUUID() error {
	lastUUID, err := u.storage.LastUUID()
	if err != nil {
		return err
	}
	u.lastUUID = lastUUID
	return nil
}
