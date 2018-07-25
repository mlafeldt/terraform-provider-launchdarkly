package launchdarkly

type Config struct {
	SDKKey string
}

type Client struct {
	config *Config
}

func (c *Config) Client() (interface{}, error) {
	return &Client{c}, nil
}
