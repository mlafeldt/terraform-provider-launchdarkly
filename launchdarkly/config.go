package launchdarkly

type Config struct {
	Token string
}

type Client struct {
	config *Config
}

func (c *Config) Client() (interface{}, error) {
	return &Client{c}, nil
}
