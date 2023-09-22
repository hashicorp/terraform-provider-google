// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package certificatemanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceGoogleCertificateManagerCertificateMap_basic(t *testing.T) {
	t.Parallel()

	project := envvar.GetTestProjectFromEnv()

	description := "My acceptance data source test certificate map"
	name := fmt.Sprintf("tf-test-certificate-map-%d", acctest.RandInt(t))
	id := fmt.Sprintf("projects/%s/locations/global/certificateMaps/%s", project, name)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleCertificateManagerCertificateMap_basic(name, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_certificate_manager_certificate_map.cert_map_data", "id", id),
					resource.TestCheckResourceAttr("data.google_certificate_manager_certificate_map.cert_map_data", "description", description),
					resource.TestCheckResourceAttr("data.google_certificate_manager_certificate_map.cert_map_data", "name", name),
				),
			},
		},
	})
}

func testAccDataSourceGoogleCertificateManagerCertificateMap_basic(certificateMapName, certificateMapDescription string) string {
	return fmt.Sprintf(`
resource "google_certificate_manager_certificate_map" "cert_map" {
	name        = "%s"
	description = "%s"
	labels      = {
		"terraform" : true,
		"acc-test"  : true,
	}
}
data "google_certificate_manager_certificate_map" "cert_map_data" {
	name = google_certificate_manager_certificate_map.cert_map.name
}
`, certificateMapName, certificateMapDescription)
}

func TestAccDataSourceGoogleCertificateManagerCertificateMap_certificateMapEntryUsingMapDatasource(t *testing.T) {
	t.Parallel()

	project := envvar.GetTestProjectFromEnv()

	certName := fmt.Sprintf("tf-test-certificate-%d", acctest.RandInt(t))
	mapEntryName := fmt.Sprintf("tf-test-certificate-map-entry-%d", acctest.RandInt(t))
	mapName := fmt.Sprintf("tf-test-certificate-map-%d", acctest.RandInt(t))
	id := fmt.Sprintf("projects/%s/locations/global/certificateMaps/%s", project, mapName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleCertificateManagerCertificateMap_certificateMapEntryUsingMapDatasource(mapName, mapEntryName, certName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_certificate_manager_certificate_map.cert_map_data", "id", id),
					resource.TestCheckResourceAttr("data.google_certificate_manager_certificate_map.cert_map_data", "name", mapName),
					resource.TestCheckResourceAttr("google_certificate_manager_certificate_map_entry.cert_map_entry", "map", mapName), // check that the certificate map entry is referencing the data source

				),
			},
		},
	})
}

func testAccDataSourceGoogleCertificateManagerCertificateMap_certificateMapEntryUsingMapDatasource(certificateMapName, certificateMapEntryName, certificateName string) string {
	return fmt.Sprintf(`
resource "google_certificate_manager_certificate_map" "cert_map" {
	name        = "%s"
	description = "certificate map example created for testing data sources in TF"
	labels      = {
		"terraform" : true,
		"acc-test"  : true,
	}
}
data "google_certificate_manager_certificate_map" "cert_map_data" {
	name = google_certificate_manager_certificate_map.cert_map.name
}
resource "google_certificate_manager_certificate" "certificate" {
	name        = "%s"
	description = "Global cert"
	self_managed {
	  pem_certificate = file("test-fixtures/cert.pem")
	  pem_private_key = file("test-fixtures/private-key.pem")
	}
}
resource "google_certificate_manager_certificate_map_entry" "cert_map_entry" {
	name        = "%s"
	description = "certificate map entry that reference a data source of certificate map and a self managed certificate"
	map = data.google_certificate_manager_certificate_map.cert_map_data.name
	certificates = [google_certificate_manager_certificate.certificate.id]
	matcher = "PRIMARY"
}
`, certificateMapName, certificateName, certificateMapEntryName)
}
