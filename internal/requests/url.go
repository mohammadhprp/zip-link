package requests

import (
	"errors"
	"net/url"
	"time"
)

type StoreURLRequest struct {
	URL       string     `bson:"url" json:"url"`
	ExpiresAt *time.Time `bson:"expires_at,omitempty" json:"expires_at"`
}

// Validate validates the StoreURLRequest struct.
func (r *StoreURLRequest) Validate() error {
	// Check if URL is not empty and is a valid URL.
	if r.URL == "" {
		return errors.New("url field is required")
	}

	if _, err := url.ParseRequestURI(r.URL); err != nil {
		return errors.New("url is not a valid format")
	}

	// Check if ExpiresAt is in the future if provided.
	if r.ExpiresAt != nil && r.ExpiresAt.Before(time.Now()) {
		return errors.New("expires_at must be a future timestamp")
	}

	return nil
}
