package website

type ContentRetrieverMonolithAdapter struct {
	Enabled bool `env:"ENABLED"`
}

type ContentRetrieverHTTPAdapter struct {
	Enabled bool   `env:"ENABLED"`
	URL     string `env:"URL"`
}

type ContentRetrieverConfig struct {
	Monolith ContentRetrieverMonolithAdapter `env:"MONOLITH"`
	HTTP     ContentRetrieverHTTPAdapter     `env:"HTTP"`
}

type Config struct {
	ContentRetriever ContentRetrieverConfig `env:"CONTENT_RETRIEVER"`
}
