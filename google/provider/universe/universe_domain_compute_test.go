// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package universe_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccUniverseDomainDisk(t *testing.T) {
	// Skip this test in all env since this can only run in specific test project.
	t.Skip()

	universeDomain := envvar.GetTestUniverseDomainFromEnv(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeDiskDestroyProducer(t),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccUniverseDomain_basic_disk(universeDomain),
			},
		},
	})
}

func TestAccDefaultUniverseDomainDisk(t *testing.T) {
	universeDomain := "googleapis.com"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeDiskDestroyProducer(t),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccUniverseDomain_basic_disk(universeDomain),
			},
		},
	})
}

func TestAccDefaultUniverseDomain_doesNotMatchExplicit(t *testing.T) {
	universeDomainFake := "fakedomain.test"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config:      testAccUniverseDomain_basic_disk(universeDomainFake),
				ExpectError: regexp.MustCompile("Universe domain mismatch"),
			},
		},
	})
}

func testAccUniverseDomain_basic_disk(universeDomain string) string {
	return fmt.Sprintf(`
provider "google" {
  universe_domain = "%s"
}
	  
resource "google_compute_instance_template" "instance_template" {
  name = "demo-this"
  machine_type = "n1-standard-1"

// boot disk
  disk {
	disk_size_gb = 20
  }

  network_interface {
	network = "default"
  }
}
`, universeDomain)
}

func testAccCheckComputeDiskDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_compute_disk" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{ComputeBasePath}}projects/{{project}}/zones/{{zone}}/disks/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("ComputeDisk still exists at %s", url)
			}
		}

		return nil
	}
}
