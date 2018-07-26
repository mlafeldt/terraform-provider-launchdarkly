package launchdarkly

import (
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/mlafeldt/go-launchdarkly/client"
)

type Config struct {
	Token string
}

type Meta struct {
	LaunchDarkly *client.Launchdarkly
	AuthInfo     runtime.ClientAuthInfoWriter
}

func (c *Config) Client() (interface{}, error) {
	authInfo := runtime.ClientAuthInfoWriterFunc(func(r runtime.ClientRequest, _ strfmt.Registry) error {
		r.SetHeaderParam("User-Agent", "Terraform")
		r.SetHeaderParam("Authorization", c.Token)
		return nil
	})

	return &Meta{
		LaunchDarkly: client.NewHTTPClient(strfmt.Default),
		AuthInfo:     authInfo,
	}, nil
}
