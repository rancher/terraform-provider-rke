package rke

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"rke": testAccProvider,
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

func hasEnvValue(envKey string) bool {
	if v := os.Getenv(envKey); v == "" {
		return false
	}
	return true
}

func testAccPreCheckEnvs(t *testing.T, keys ...string) {
	for _, env := range keys {
		if !hasEnvValue(env) {
			t.Fatal(fmt.Sprintf("%s must be set for acceptance tests", env))
		}
	}
}

func testAccPreCheck(t *testing.T) {
	testAccPreCheckEnvs(t, envRKENodeAddr, envRKENodeUser, envRKENodeSSHKey)
	rkeLogBuf.Reset()
}

func testAccPreCheckForMultiNodes(t *testing.T) {
	var envKeys []string
	baseKeys := []string{envRKENodeAddr, envRKENodeUser, envRKENodeSSHKey}
	for i := 0; i < 2; i++ {
		for _, key := range baseKeys {
			envKeys = append(envKeys, fmt.Sprintf("%s_%d", key, i))
		}
	}
	testAccPreCheckEnvs(t, envKeys...)
}
