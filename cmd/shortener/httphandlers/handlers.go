package httphandlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"

	"github.com/MrTomSawyer/url-shortener/cmd/shortener/tools"
)

type URLData struct {
	URL string
}

var vault = make(map[string]string)

func HTTPHandler(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		http.Error(res, "Error parsing params", http.StatusInternalServerError)
	}
	switch req.Method {
	case "POST":
		ShortenUrl(res, req)
	case "GET":
		fmt.Println("GET")
		GetOriginalURL(res, req)
	}
}

func ShortenUrl(res http.ResponseWriter, req *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("Error closing Body: %v", err)
		}
	}(req.Body)

	if req.Method != http.MethodPost && req.Method != http.MethodGet {
		http.Error(res, "Only POST and GET methods allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var data URLData
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(res, "Error unmarshalling request body", http.StatusBadRequest)
		return
	}

	shortUrl, hash := tools.Shorten(data.URL)
	if _, ok := vault[hash]; !ok {
		vault[hash] = shortUrl
	}
	res.WriteHeader(http.StatusCreated)
	if _, err = res.Write([]byte(shortUrl)); err != nil {
		http.Error(res, "Error sending response", http.StatusInternalServerError)
	}
}

func GetOriginalURL(res http.ResponseWriter, req *http.Request) {
	_, id := path.Split(req.URL.Path)
	if value, ok := vault[id]; ok {
		res.Header().Set("Location", value)
		res.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		http.Error(res, "URL not found", http.StatusNotFound)
	}
}
