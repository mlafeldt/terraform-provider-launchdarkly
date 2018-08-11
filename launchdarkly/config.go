package launchdarkly

import (
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mlafeldt/go-launchdarkly/client"
)

// Config is used to configure the creation of a LaunchDarkly client.
type Config struct {
	Token string
}

// Meta is used by the provider to access the LaunchDarkly API.
type Meta struct {
	LaunchDarkly *client.Launchdarkly
	AuthInfo     runtime.ClientAuthInfoWriter
}

// Client returns a configured LaunchDarkly client.
func (c *Config) Client() (interface{}, error) {
	authInfo := runtime.ClientAuthInfoWriterFunc(func(r runtime.ClientRequest, _ strfmt.Registry) error {
		r.SetHeaderParam("User-Agent", "Terraform/"+terraform.VersionString())
		r.SetHeaderParam("Authorization", c.Token)
		return nil
	})

	return &Meta{
		LaunchDarkly: client.NewHTTPClient(strfmt.Default),
		AuthInfo:     authInfo,
	}, nil
}
