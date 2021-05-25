package truenasapi

/* Config (Endpoint, api key, etc) */
type Config struct {
	Host   string
	APIKey string
}

func NewDefaultConfig() (*Config, error) {

	return &Config{
		Host:   "10.50.1.20",
		APIKey: "3-YKjefi3YT22oOfpVdOXaognmHxvA60xSOoP0Dj1fbJ3qdJCPs1nZtqbgDhtT4Wgr",
	}, nil
}
