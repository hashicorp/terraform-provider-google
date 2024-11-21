// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package acctest

import (
	"context"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/provider"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var TestAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	configs = make(map[string]*transport_tpg.Config)
	sources = make(map[string]VcrSource)
	testAccProvider = provider.Provider()
	TestAccProviders = map[string]*schema.Provider{
		"google": testAccProvider,
	}
}

// GoogleProviderConfig returns a configured SDKv2 provider.
// This function is typically used in CheckDestroy functions in acceptance tests. The provider client is used to make GET requests to check a resource is destroyed.
// Either a preexisting configured SDKv2 provider for the given test name is returned, or a new one is configured with empty (but non-nil) terraform.ResourceConfig
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

// AccTestPreCheck_AdcCredentialsOnly is a PreCheck function for acceptance tests that use ADCs when
func AccTestPreCheck_AdcCredentialsOnly(t *testing.T) {
	if v := os.Getenv("GOOGLE_CREDENTIALS_FILE"); v != "" {
		t.Log("Ignoring GOOGLE_CREDENTIALS_FILE; acceptance test doesn't use credentials other than ADCs")
	}

	// Fail on set creds
	if v := transport_tpg.MultiEnvSearch(envvar.CredsEnvVarsExcludingAdcs()); v != "" {
		t.Fatalf("This acceptance test only uses ADCs, so all of %s must be unset", strings.Join(envvar.CredsEnvVarsExcludingAdcs(), ", "))
	}

	// Fail on ADC ENV not set
	if v := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); v == "" {
		t.Fatalf("GOOGLE_APPLICATION_CREDENTIALS must be set for acceptance tests that are dependent on ADCs")
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
