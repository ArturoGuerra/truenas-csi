package truenasapi

import "os"

/* Config (Endpoint, api key, etc) */
type Config struct {
	Host   string
	APIKey string
}

func NewDefaultConfig() (*Config, error) {
	return &Config{
		Host:   os.Getenv("TRUENAS_HOST"),
		APIKey: os.Getenv("TRUENAS_API_KEY"),
	}, nil
}
