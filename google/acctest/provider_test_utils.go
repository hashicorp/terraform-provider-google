// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package acctest

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/provider"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var TestAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	configs = make(map[string]*transport_tpg.Config)
	fwProviders = make(map[string]*frameworkTestProvider)
	sources = make(map[string]VcrSource)
	testAccProvider = provider.Provider()
	TestAccProviders = map[string]*schema.Provider{
		"google": testAccProvider,
	}
}

func GoogleProviderConfig(t *testing.T) *transport_tpg.Config {
	configsLock.RLock()
	config, ok := configs[t.Name()]
	configsLock.RUnlock()
	if ok {
		return config
	}

	sdkProvider := provider.Provider()
	rc := terraform.ResourceConfig{}
	sdkProvider.Configure(context.Background(), &rc)
	return sdkProvider.Meta().(*transport_tpg.Config)
}

func AccTestPreCheck(t *testing.T) {
	if v := os.Getenv("GOOGLE_CREDENTIALS_FILE"); v != "" {
		creds, err := ioutil.ReadFile(v)
		if err != nil {
			t.Fatalf("Error reading GOOGLE_CREDENTIALS_FILE path: %s", err)
		}
		os.Setenv("GOOGLE_CREDENTIALS", string(creds))
	}

	if v := transport_tpg.MultiEnvSearch(envvar.CredsEnvVars); v == "" {
		t.Fatalf("One of %s must be set for acceptance tests", strings.Join(envvar.CredsEnvVars, ", "))
	}

	if v := transport_tpg.MultiEnvSearch(envvar.ProjectEnvVars); v == "" {
		t.Fatalf("One of %s must be set for acceptance tests", strings.Join(envvar.ProjectEnvVars, ", "))
	}

	if v := transport_tpg.MultiEnvSearch(envvar.RegionEnvVars); v == "" {
		t.Fatalf("One of %s must be set for acceptance tests", strings.Join(envvar.RegionEnvVars, ", "))
	}

	if v := transport_tpg.MultiEnvSearch(envvar.ZoneEnvVars); v == "" {
		t.Fatalf("One of %s must be set for acceptance tests", strings.Join(envvar.ZoneEnvVars, ", "))
	}
}

// GetTestRegion has the same logic as the provider's GetRegion, to be used in tests.
func GetTestRegion(is *terraform.InstanceState, config *transport_tpg.Config) (string, error) {
	if res, ok := is.Attributes["region"]; ok {
		return res, nil
	}
	if config.Region != "" {
		return config.Region, nil
	}
	return "", fmt.Errorf("%q: required field is not set", "region")
}

// GetTestProject has the same logic as the provider's GetProject, to be used in tests.
func GetTestProject(is *terraform.InstanceState, config *transport_tpg.Config) (string, error) {
	if res, ok := is.Attributes["project"]; ok {
		return res, nil
	}
	if config.Project != "" {
		return config.Project, nil
	}
	return "", fmt.Errorf("%q: required field is not set", "project")
}

// Some tests fail during VCR. One common case is race conditions when creating resources.
// If a test config adds two fine-grained resources with the same parent it is undefined
// which will be created first, causing VCR to fail ~50% of the time
func SkipIfVcr(t *testing.T) {
	if IsVcrEnabled() {
		t.Skipf("VCR enabled, skipping test: %s", t.Name())
	}
}

func SleepInSecondsForTest(t int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		time.Sleep(time.Duration(t) * time.Second)
		return nil
	}
}
