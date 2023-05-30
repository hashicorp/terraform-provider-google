// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Tests schema version migration by creating a certificate with an old version of the provider (4.59.0)
// and then updating it with the current version the provider.
func TestAccCertificateManagerCertificate_migration(t *testing.T) {
	SkipIfVcr(t)
	t.Parallel()
	name := fmt.Sprintf("tf-test-%d", RandInt(t))

	oldVersion := map[string]resource.ExternalProvider{
		"google": {
			VersionConstraint: "4.59.0", // a version that doesn't support location yet.
			Source:            "registry.terraform.io/hashicorp/google",
		},
	}
	newVersion := map[string]func() (*schema.Provider, error){
		"mynewprovider": func() (*schema.Provider, error) { return TestAccProviders["google"], nil },
	}

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { AccTestPreCheck(t) },
		CheckDestroy: testAccCheckCertificateManagerCertificateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:            configWithDescritption(name),
				ExternalProviders: oldVersion,
			},
			{
				ResourceName:            "google_certificate_manager_certificate.default",
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"location", "self_managed"},
				ExternalProviders:       oldVersion,
			},
			{
				Config:            newConfigWithDescription(name),
				ProviderFactories: newVersion,
			},
			{
				ResourceName:            "google_certificate_manager_certificate.default",
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"location", "self_managed"},
				ProviderFactories:       newVersion,
			},
		},
	})
}

func configWithDescritption(name string) string {
	return fmt.Sprintf(`
	resource "google_certificate_manager_certificate" "default" {
		name        = "%s"
		description = "Global cert"
		self_managed {
		  pem_certificate = file("test-fixtures/certificatemanager/cert.pem")
		  pem_private_key = file("test-fixtures/certificatemanager/private-key.pem")
		}
	}
	`, name)
}

func newConfigWithDescription(name string) string {
	return fmt.Sprintf(`
	provider "mynewprovider" {}
	
	resource "google_certificate_manager_certificate" "default" {
		provider    = mynewprovider
		name        = "%s"
		description = "Migrated Global cert"
		self_managed {
		  pem_certificate = file("test-fixtures/certificatemanager/cert.pem")
		  pem_private_key = file("test-fixtures/certificatemanager/private-key.pem")
		}
	}
	`, name)
}
