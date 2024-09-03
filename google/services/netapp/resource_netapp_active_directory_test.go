// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package netapp_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccNetappActiveDirectory_activeDirectory_FullUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetappActiveDirectory_activeDirectoryCreateExample_Full(context),
			},
			{
				ResourceName:            "google_netapp_active_directory.test_active_directory_full",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "pass", "labels", "terraform_labels"},
			},
			{
				Config: testAccNetappActiveDirectory_activeDirectoryCreateExample_Update(context),
			},
			{
				ResourceName:            "google_netapp_active_directory.test_active_directory_full",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "pass", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccNetappActiveDirectory_activeDirectoryCreateExample_Full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_netapp_active_directory" "test_active_directory_full" {
    name = "tf-test-test-active-directory-full%{random_suffix}"
    location = "us-central1"
    domain = "ad.internal"
    dns = "172.30.64.3"
    net_bios_prefix = "smbserver"
    username = "user"
    password = "pass"
    aes_encryption         = false
    backup_operators       = ["test1", "test2"]
    administrators         = ["test1", "test2"]
    description            = "ActiveDirectory is the public representation of the active directory config."
    encrypt_dc_connections = false
    kdc_hostname           = "hostname"
    kdc_ip                 = "10.10.0.11"
    labels                 = { 
        "foo": "bar"
    }
    ldap_signing           = false
    nfs_users_with_ldap    = false
    organizational_unit    = "CN=Computers"
    security_operators     = ["test1", "test2"]
    site                   = "test-site"
  }
`, context)
}

func testAccNetappActiveDirectory_activeDirectoryCreateExample_Update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_netapp_active_directory" "test_active_directory_full" {
    name = "tf-test-test-active-directory-full%{random_suffix}"
    location = "us-central1"
    domain = "ad.internal"
    dns = "172.30.64.3"
    net_bios_prefix = "smbup"
    username = "user"
    password = "pass"
    aes_encryption         = false
    backup_operators       = ["test1", "test2"]
    administrators         = ["test1", "test2"]
    description            = "ActiveDirectory is the public representation of the active directory config."
    encrypt_dc_connections = false
    kdc_hostname           = "hostname"
    kdc_ip                 = "10.10.0.11"
    labels                 = { 
        "foo": "bar"
    }
    ldap_signing           = true
    nfs_users_with_ldap    = true
    organizational_unit    = "CN=Computers"
    security_operators     = ["test1", "test2"]
    site                   = "test-site"
  }
`, context)
}
