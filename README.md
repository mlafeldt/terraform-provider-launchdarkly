# Terraform Provider for LaunchDarkly

[![CircleCI](https://circleci.com/gh/mlafeldt/terraform-provider-launchdarkly.svg?style=svg)](https://circleci.com/gh/mlafeldt/terraform-provider-launchdarkly)

## Using the Provider

### Supported Resources

- `launchdarkly_project`
- `launchdarkly_environment`
- `launchdarkly_feature_flag`

### Example

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

## Installing the Provider

You'll first need [Go](http://www.golang.org) installed on your machine (version 1.9+ is *required*). You'll also need to correctly set up a [GOPATH](http://golang.org/doc/code.html#GOPATH).

To build and install the provider from source, do the following:

```console
git clone https://github.com/mlafeldt/terraform-provider-launchdarkly
cd terraform-provider-launchdarkly/
make install
```

## Developing the Provider

To compile the provider, run `make build`. This will build the provider and put the provider binary in the current directory.

```console
make build
...
./terraform-provider-launchdarkly
...
```

In order to test the binary with Terraform, create a .tf file in the current directory and run `make apply` or `make destroy`.

```console
vim example.tf
...
make apply
...
make destroy
...
```

## Known Issues

- Creating a new project via the LaunchDarkly API will, unfortunately, also create two environments by default (Production and Test). It's not possible to delete both since at least one must exist. The `launchdarkly_project` and `launchdarkly_environment` resources don't handle this case in a good way yet. As a workaround, it's possible to import these special environments via `terraform import`. See <https://github.com/launchdarkly/ld-openapi/issues/7>.
- Even if we could solve the previous issue, there would still be the problem that there must be at least one project, making it impossible to have a clean slate using Terraform.
- `launchdarkly_feature_flag` doesn't support custom variations yet (bool only).
