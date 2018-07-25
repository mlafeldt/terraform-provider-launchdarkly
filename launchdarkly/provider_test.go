package launchdarkly

import (
	"errors"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"launchdarkly": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("LAUNCHDARKLY_SDK_KEY"); v == "" {
		t.Fatal("LAUNCHDARKLY_SDK_KEY must be set for acceptance tests")
	}
}

func testAccAssert(msg string, f func() bool) resource.TestCheckFunc {
	return func(*terraform.State) error {
		if f() {
			return nil
		}

		return errors.New("assertion failed: " + msg)
	}
}
