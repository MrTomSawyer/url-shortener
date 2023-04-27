package httphandlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/MrTomSawyer/url-shortener/cmd/shortener/tools"
)

type URLData struct {
	URL string
}

var vault map[string]string = make(map[string]string)

func ShortenUrl(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	if req.Method != http.MethodPost {
		http.Error(res, "Only POST method allowed", http.StatusMethodNotAllowed)
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
		http.Error(res, "Error unmarshaling request body", http.StatusBadRequest)
		return
	}

	shortUrl, hash := tools.Shorten(data.URL)
	vault[hash] = shortUrl
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(shortUrl))
}

func GetOriginalURL(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		http.Error(res, "Failed to parse request params", http.StatusInternalServerError)
		return
	}

	id := req.FormValue("id")
	if value, ok := vault[id]; ok {
		res.Header().Set("Location", value)
		res.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		http.Error(res, "URL not found", http.StatusNotFound)
	}
}
