# Terraform Provider for LaunchDarkly

[![CircleCI](https://circleci.com/gh/mlafeldt/terraform-provider-launchdarkly.svg?style=svg)](https://circleci.com/gh/mlafeldt/terraform-provider-launchdarkly)

## Supported Resources

- `launchdarkly_project`
- `launchdarkly_environment`
- `launchdarkly_feature_flag`

## Using the Provider

```hcl
resource "launchdarkly_project" "example" {
  name = "Example Project"
  key  = "example"
}

resource "launchdarkly_environment" "staging" {
  name  = "Staging"
  key   = "staging"
  color = "0000ff"

  project = "${launchdarkly_project.example.key}"
}

resource launchdarkly_feature_flag "example" {
  key  = "some.example.flag"
  name = "Some example flag"

  tags = [
    "hello",
    "world",
  ]

  project = "${launchdarkly_project.example.key}"
}

output "staging_api_key" {
  value = "${launchdarkly_environment.staging.api_key}"
}
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.9+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the current directory.

```console
$ make build
...
$ ./terraform-provider-launchdarkly
...
```

In order to test the provider, you can simply run `make test`.

```console
$ make test
```

In order to run the full suite of acceptance tests, run `make testacc`.

```console
$ make testacc
```
