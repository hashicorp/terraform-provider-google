// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

// Deprecated: For backward compatibility CheckDataSourceStateMatchesResourceState is still working,
// but all new code should use CheckDataSourceStateMatchesResourceState in the acctest package instead.
func CheckDataSourceStateMatchesResourceState(dataSourceName, resourceName string) func(*terraform.State) error {
	return acctest.CheckDataSourceStateMatchesResourceState(dataSourceName, resourceName)
}

// Deprecated: For backward compatibility CheckDataSourceStateMatchesResourceStateWithIgnores is still working,
// but all new code should use CheckDataSourceStateMatchesResourceStateWithIgnores in the acctest package instead.
func CheckDataSourceStateMatchesResourceStateWithIgnores(dataSourceName, resourceName string, ignoreFields map[string]struct{}) func(*terraform.State) error {
	return acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(dataSourceName, resourceName, ignoreFields)
}

// General test utils

func RandString(t *testing.T, length int) string {
	return acctest.RandString(t, length)
}

func RandInt(t *testing.T) int {
	return acctest.RandInt(t)
}

// ProtoV5ProviderFactories returns a muxed ProviderServer that uses the provider code from this repo (SDK and plugin-framework).
// Used to set ProtoV5ProviderFactories in a resource.TestStep within an acceptance test.
func ProtoV5ProviderFactories(t *testing.T) map[string]func() (tfprotov5.ProviderServer, error) {
	return acctest.ProtoV5ProviderFactories(t)
}

// ProtoV5ProviderBetaFactories returns the same as ProtoV5ProviderFactories only the provider is mapped with
// "google-beta" to ensure that registry examples use `google-beta` if the example is versioned as beta;
// normal beta tests should continue to use ProtoV5ProviderFactories
func ProtoV5ProviderBetaFactories(t *testing.T) map[string]func() (tfprotov5.ProviderServer, error) {
	return acctest.ProtoV5ProviderBetaFactories(t)
}

func serviceAccountCanonicalEmail(account string) string {
	return envvar.ServiceAccountCanonicalEmail(account)
}

func getResourceAttributes(n string, s *terraform.State) (map[string]string, error) {
	return tpgresource.GetResourceAttributes(n, s)
}

// Deprecated: For backward compatibility testBucketName is still working,
// but all new code should use TestBucketName in the acctest package instead.
func testBucketName(t *testing.T) string {
	return acctest.TestBucketName(t)
}

// Deprecated: For backward compatibility createZIPArchiveForCloudFunctionSource is still working,
// but all new code should use CreateZIPArchiveForCloudFunctionSource in the acctest package instead.
func createZIPArchiveForCloudFunctionSource(t *testing.T, sourcePath string) string {
	return acctest.CreateZIPArchiveForCloudFunctionSource(t, sourcePath)
}
