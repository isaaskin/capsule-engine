package registryhandler

import (
	"time"
)

type Repository struct {
	Name           string    `json:"name"`
	Namespace      string    `json:"namespace"`
	RepositoryType string    `json:"repository_type"`
	Status         int       `json:"status"`
	StatusDesc     string    `json:"status_description"`
	Description    string    `json:"description"`
	IsPrivate      bool      `json:"is_private"`
	StarCount      int       `json:"star_count"`
	PullCount      int       `json:"pull_count"`
	LastUpdated    time.Time `json:"last_updated"`
	DateRegistered time.Time `json:"date_registered"`
	Affiliation    string    `json:"affiliation"`
	MediaTypes     []string  `json:"media_types"`
	ContentTypes   []string  `json:"content_types"`
}

type RepositoryList struct {
	Count    int          `json:"count"`
	Next     interface{}  `json:"next"`     // May be a URL or null, hence `interface{}`
	Previous interface{}  `json:"previous"` // May be a URL or null, hence `interface{}`
	Results  []Repository `json:"results"`
}
