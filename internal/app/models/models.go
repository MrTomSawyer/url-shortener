// Package models defines various data structures used for representing request and response data in the application.
package models

// ShortenRequest represents a request to shorten a URL.
type ShortenRequest struct {
	URL string `json:"url"`
}

// ShortenResponse represents the response containing the shortened URL.
type ShortenResponse struct {
	Result string `json:"result"`
}

// URLJson represents a JSON representation of a URL.
type URLJson struct {
	UUID        int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// BatchURLRequest represents a request to shorten a URL in a batch.
type BatchURLRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// BatchURLResponse represents the response containing the shortened URL for a batch request.
type BatchURLResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// TempURLBatchRequest represents a temporary URL batch request.
type TempURLBatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
	OriginalURL   string `json:"original_url"`
}

// URLJsonResponse represents a JSON representation of a URL in a response.
type URLJsonResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// UserURLs represents a user's URLs.
type UserURLs struct {
	UserID string   `json:"user_id"`
	URLs   []string `json:"urls"`
}
