// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package fwprovider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

// TestAccFwProvider_region is a series of acc tests asserting how the plugin-framework provider handles credentials arguments
// It is PF specific because the HCL used uses a PF-implemented data source
// It is a counterpart to TestAccSdkProvider_region
func TestAccFwProvider_region(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		"config takes precedence over environment variables":                          testAccFwProvider_region_configPrecedenceOverEnvironmentVariables,
		"when region is unset in the config, environment variables are used":          testAccFwProvider_region_precedenceOrderEnvironmentVariables,
		"when region is set to an empty string in the config the value isn't ignored": testAccFwProvider_region_emptyStringValidation,
		"region values can be supplied as a self link, but are transformed":           testAccFwProvider_region_selfLinks,
	}

	for name, tc := range testCases {
		// shadow the tc variable into scope so that when
		// the loop continues, if t.Run hasn't executed tc(t)
		// yet, we don't have a race condition
		// see https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc(t)
		})
	}
}

// TestAccFwProvider_zone is a series of acc tests asserting how the plugin-framework provider handles credentials arguments
// It is PF specific because the HCL used uses a PF-implemented data source
// It is a counterpart to TestAccSdkProvider_zone
func TestAccFwProvider_zone(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		"config takes precedence over environment variables":                        testAccFwProvider_zone_configPrecedenceOverEnvironmentVariables,
		"when zone is unset in the config, environment variables are used":          testAccFwProvider_zone_precedenceOrderEnvironmentVariables,
		"when zone is set to an empty string in the config the value isn't ignored": testAccFwProvider_zone_emptyStringValidation,

		// different behavior vs region
		"zone values can be supplied as a self link and are NOT transformed": testAccFwProvider_zone_selfLinks,
	}

	for name, tc := range testCases {
		// shadow the tc variable into scope so that when
		// the loop continues, if t.Run hasn't executed tc(t)
		// yet, we don't have a race condition
		// see https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc(t)
		})
	}
}

func testAccFwProvider_region_configPrecedenceOverEnvironmentVariables(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	region := envvar.GetTestRegionFromEnv()

	// ensure all possible region env vars set; show they aren't used
	for _, v := range envvar.RegionEnvVars {
		t.Setenv(v, region)
	}

	providerRegion := "foobar"

	context := map[string]interface{}{
		"region": providerRegion,
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Apply-time error; bad value in config is used over of good values in ENVs
				Config: testAccFwProvider_regionInProviderBlock(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_plugin_framework.default", "region", providerRegion),
				),
			},
		},
	})
}

func testAccFwProvider_zone_configPrecedenceOverEnvironmentVariables(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	zone := envvar.GetTestZoneFromEnv()

	// ensure all possible zone env vars set; show they aren't used
	for _, v := range envvar.ZoneEnvVars {
		t.Setenv(v, zone)
	}

	providerZone := "foobar"

	context := map[string]interface{}{
		"zone": providerZone,
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Apply-time error; bad value in config is used over of good values in ENVs
				Config: testAccFwProvider_zoneInProviderBlock(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_plugin_framework.default", "zone", providerZone),
				),
			},
		},
	})
}

func testAccFwProvider_region_precedenceOrderEnvironmentVariables(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API
	/*
		These are all the ENVs for region, and they are in order of precedence.
		GOOGLE_REGION
		GCLOUD_REGION
		CLOUDSDK_COMPUTE_REGION
	*/

	GOOGLE_REGION := "GOOGLE_REGION"
	GCLOUD_REGION := "GCLOUD_REGION"
	CLOUDSDK_COMPUTE_REGION := "CLOUDSDK_COMPUTE_REGION"

	context := map[string]interface{}{}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// GOOGLE_REGION is used 1st if set
				PreConfig: func() {
					t.Setenv("GOOGLE_REGION", GOOGLE_REGION) //used
					t.Setenv("GCLOUD_REGION", GCLOUD_REGION)
					t.Setenv("CLOUDSDK_COMPUTE_REGION", CLOUDSDK_COMPUTE_REGION)
				},
				Config: testAccFwProvider_regionInEnvsOnly(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_plugin_framework.default", "region", GOOGLE_REGION),
				),
			},
			{
				// GCLOUD_REGION is used 2nd
				PreConfig: func() {
					// unset
					t.Setenv("GOOGLE_REGION", "")
					// set
					t.Setenv("GCLOUD_REGION", GCLOUD_REGION) //used
					t.Setenv("CLOUDSDK_COMPUTE_REGION", CLOUDSDK_COMPUTE_REGION)
				},
				Config: testAccFwProvider_regionInEnvsOnly(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_plugin_framework.default", "region", GCLOUD_REGION),
				),
			},
			{
				// GOOGLE_CLOUD_KEYFILE_JSON is used 3rd
				PreConfig: func() {
					// unset
					t.Setenv("GOOGLE_REGION", "")
					t.Setenv("GCLOUD_REGION", "")
					// set
					t.Setenv("CLOUDSDK_COMPUTE_REGION", CLOUDSDK_COMPUTE_REGION) //used
				},
				Config: testAccFwProvider_regionInEnvsOnly(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_plugin_framework.default", "region", CLOUDSDK_COMPUTE_REGION),
				),
			},
		},
	})
}

func testAccFwProvider_zone_precedenceOrderEnvironmentVariables(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API
	/*
		These are all the ENVs for zone, and they are in order of precedence.
		GOOGLE_ZONE
		GCLOUD_ZONE
		CLOUDSDK_COMPUTE_ZONE
	*/

	GOOGLE_ZONE := "GOOGLE_ZONE"
	GCLOUD_ZONE := "GCLOUD_ZONE"
	CLOUDSDK_COMPUTE_ZONE := "CLOUDSDK_COMPUTE_ZONE"

	context := map[string]interface{}{}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// GOOGLE_ZONE is used 1st if set
				PreConfig: func() {
					t.Setenv("GOOGLE_ZONE", GOOGLE_ZONE) //used
					t.Setenv("GCLOUD_ZONE", GCLOUD_ZONE)
					t.Setenv("CLOUDSDK_COMPUTE_ZONE", CLOUDSDK_COMPUTE_ZONE)
				},
				Config: testAccFwProvider_zoneInEnvsOnly(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_plugin_framework.default", "zone", GOOGLE_ZONE),
				),
			},
			{
				// GCLOUD_ZONE is used 2nd
				PreConfig: func() {
					// unset
					t.Setenv("GOOGLE_ZONE", "")
					// set
					t.Setenv("GCLOUD_ZONE", GCLOUD_ZONE) //used
					t.Setenv("CLOUDSDK_COMPUTE_ZONE", CLOUDSDK_COMPUTE_ZONE)
				},
				Config: testAccFwProvider_zoneInEnvsOnly(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_plugin_framework.default", "zone", GCLOUD_ZONE),
				),
			},
			{
				// GOOGLE_CLOUD_KEYFILE_JSON is used 3rd
				PreConfig: func() {
					// unset
					t.Setenv("GOOGLE_ZONE", "")
					t.Setenv("GCLOUD_ZONE", "")
					// set
					t.Setenv("CLOUDSDK_COMPUTE_ZONE", CLOUDSDK_COMPUTE_ZONE) //used
				},
				Config: testAccFwProvider_zoneInEnvsOnly(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_plugin_framework.default", "zone", CLOUDSDK_COMPUTE_ZONE),
				),
			},
		},
	})
}

func testAccFwProvider_region_emptyStringValidation(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	context := map[string]interface{}{
		"region": "",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccFwProvider_regionInProviderBlock(context),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("expected a non-empty string"),
			},
		},
	})
}

func testAccFwProvider_zone_emptyStringValidation(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	context := map[string]interface{}{
		"zone": "",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccFwProvider_zoneInProviderBlock(context),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("expected a non-empty string"),
			},
		},
	})
}

func testAccFwProvider_region_selfLinks(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	selfLink := "https://www.googleapis.com/compute/v1/projects/my-project/regions/us-central1"
	region := "us-central1"

	context := map[string]interface{}{
		"region": selfLink,
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFwProvider_regionInProviderBlock(context),
				Check: resource.ComposeTestCheckFunc(
					// output value is transformed
					resource.TestCheckResourceAttr("data.google_provider_config_plugin_framework.default", "region", region),
				),
			},
		},
	})
}

func testAccFwProvider_zone_selfLinks(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	selfLink := "https://www.googleapis.com/compute/v1/projects/my-project/zones/us-central1-a"

	context := map[string]interface{}{
		"zone": selfLink,
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFwProvider_zoneInProviderBlock(context),
				Check: resource.ComposeTestCheckFunc(
					// output value is transformed
					resource.TestCheckResourceAttr("data.google_provider_config_plugin_framework.default", "zone", selfLink),
				),
			},
		},
	})
}

// testAccFwProvider_regionInProviderBlock allows setting the region argument in a provider block.
// This function uses data.google_provider_config_plugin_framework because it is implemented with the plugin-framework
func testAccFwProvider_regionInProviderBlock(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
	region = "%{region}"
}

data "google_provider_config_plugin_framework" "default" {}

output "region" {
  value = data.google_provider_config_plugin_framework.default.region
  sensitive = true
}
`, context)
}

// testAccFwProvider_regionInEnvsOnly allows testing when the region argument
// is only supplied via ENVs
func testAccFwProvider_regionInEnvsOnly(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_provider_config_plugin_framework" "default" {}

output "region" {
  value = data.google_provider_config_plugin_framework.default.region
  sensitive = true
}
`, context)
}

// testAccFwProvider_zoneInProviderBlock allows setting the zone argument in a provider block.
// This function uses data.google_provider_config_plugin_framework because it is implemented with the plugin-framework
func testAccFwProvider_zoneInProviderBlock(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
	zone = "%{zone}"
}

data "google_provider_config_plugin_framework" "default" {}

output "zone" {
  value = data.google_provider_config_plugin_framework.default.zone
  sensitive = true
}
`, context)
}

// testAccFwProvider_zoneInEnvsOnly allows testing when the zone argument
// is only supplied via ENVs
func testAccFwProvider_zoneInEnvsOnly(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_provider_config_plugin_framework" "default" {}

output "zone" {
  value = data.google_provider_config_plugin_framework.default.zone
  sensitive = true
}
`, context)
}
