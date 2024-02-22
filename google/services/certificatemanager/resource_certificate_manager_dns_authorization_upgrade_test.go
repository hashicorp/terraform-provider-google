// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package certificatemanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

// Tests schema version migration by creating a dns authorization with an old version of the provider (5.15.0)
// and then updating it with the current version the provider.
func TestAccCertificateManagerDnsAuthorization_migration(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()
	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	oldVersion := map[string]resource.ExternalProvider{
		"google": {
			VersionConstraint: "5.15.0", // a version that doesn't support location yet.
			Source:            "registry.terraform.io/hashicorp/google",
		},
	}
	newVersion := map[string]func() (*schema.Provider, error){
		"mynewprovider": func() (*schema.Provider, error) { return acctest.TestAccProviders["google"], nil },
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:     func() { acctest.AccTestPreCheck(t) },
		CheckDestroy: testAccCheckCertificateManagerDnsAuthorizationDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:            dnsAuthorizationResourceConfig(name),
				ExternalProviders: oldVersion,
			},
			{
				ResourceName:            "google_certificate_manager_dns_authorization.default",
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"location"},
				ExternalProviders:       oldVersion,
			},
			{
				Config:            dnsAuthorizationResourceConfigUpdated(name),
				ProviderFactories: newVersion,
			},
			{
				ResourceName:            "google_certificate_manager_dns_authorization.default",
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"location"},
				ProviderFactories:       newVersion,
			},
		},
	})
}

func dnsAuthorizationResourceConfig(name string) string {
	return fmt.Sprintf(`
	resource "google_certificate_manager_dns_authorization" "default" {
		name        = "%s"
		description = "The default dns"
		domain      = "domain.hashicorptest.com"
	  }
	`, name)
}

func dnsAuthorizationResourceConfigUpdated(name string) string {
	return fmt.Sprintf(`
	provider "mynewprovider" {}
	
	resource "google_certificate_manager_dns_authorization" "default" {
		provider = mynewprovider
		name        = "%s"
		description = "The migrated default dns"
		domain      = "domain.hashicorptest.com"
	  }
	`, name)
}
