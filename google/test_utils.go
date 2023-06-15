// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"context"
	"math/rand"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	acctest_tpg "github.com/hashicorp/terraform-provider-google/google/acctest"
)

// Deprecated: For backward compatibility CheckDataSourceStateMatchesResourceState is still working,
// but all new code should use CheckDataSourceStateMatchesResourceState in the acctest package instead.
func CheckDataSourceStateMatchesResourceState(dataSourceName, resourceName string) func(*terraform.State) error {
	return acctest_tpg.CheckDataSourceStateMatchesResourceState(dataSourceName, resourceName)
}

// Deprecated: For backward compatibility CheckDataSourceStateMatchesResourceStateWithIgnores is still working,
// but all new code should use CheckDataSourceStateMatchesResourceStateWithIgnores in the acctest package instead.
func CheckDataSourceStateMatchesResourceStateWithIgnores(dataSourceName, resourceName string, ignoreFields map[string]struct{}) func(*terraform.State) error {
	return acctest_tpg.CheckDataSourceStateMatchesResourceStateWithIgnores(dataSourceName, resourceName, ignoreFields)
}

// General test utils

// MuxedProviders returns the correct test provider (between the sdk version or the framework version)
func MuxedProviders(testName string) (func() tfprotov5.ProviderServer, error) {
	ctx := context.Background()

	providers := []func() tfprotov5.ProviderServer{
		providerserver.NewProtocol5(NewFrameworkTestProvider(testName)), // framework provider
		GetSDKProvider(testName).GRPCProvider,                           // sdk provider
	}

	muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)

	if err != nil {
		return nil, err
	}

	return muxServer.ProviderServer, nil
}

func RandString(t *testing.T, length int) string {
	if !acctest_tpg.IsVcrEnabled() {
		return acctest.RandString(length)
	}
	envPath := os.Getenv("VCR_PATH")
	vcrMode := os.Getenv("VCR_MODE")
	s, err := vcrSource(t, envPath, vcrMode)
	if err != nil {
		// At this point we haven't created any resources, so fail fast
		t.Fatal(err)
	}

	r := rand.New(s.source)
	result := make([]byte, length)
	set := "abcdefghijklmnopqrstuvwxyz012346789"
	for i := 0; i < length; i++ {
		result[i] = set[r.Intn(len(set))]
	}
	return string(result)
}

func RandInt(t *testing.T) int {
	if !acctest_tpg.IsVcrEnabled() {
		return acctest.RandInt()
	}
	envPath := os.Getenv("VCR_PATH")
	vcrMode := os.Getenv("VCR_MODE")
	s, err := vcrSource(t, envPath, vcrMode)
	if err != nil {
		// At this point we haven't created any resources, so fail fast
		t.Fatal(err)
	}

	return rand.New(s.source).Int()
}

// ProtoV5ProviderFactories returns a muxed ProviderServer that uses the provider code from this repo (SDK and plugin-framework).
// Used to set ProtoV5ProviderFactories in a resource.TestStep within an acceptance test.
func ProtoV5ProviderFactories(t *testing.T) map[string]func() (tfprotov5.ProviderServer, error) {
	return map[string]func() (tfprotov5.ProviderServer, error){
		"google": func() (tfprotov5.ProviderServer, error) {
			provider, err := MuxedProviders(t.Name())
			return provider(), err
		},
	}
}

// ProtoV5ProviderBetaFactories returns the same as ProtoV5ProviderFactories only the provider is mapped with
// "google-beta" to ensure that registry examples use `google-beta` if the example is versioned as beta;
// normal beta tests should continue to use ProtoV5ProviderFactories
func ProtoV5ProviderBetaFactories(t *testing.T) map[string]func() (tfprotov5.ProviderServer, error) {
	return map[string]func() (tfprotov5.ProviderServer, error){}
}
