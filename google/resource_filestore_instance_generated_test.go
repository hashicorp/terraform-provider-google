// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccFilestoreInstance_filestoreInstanceBasicExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckFilestoreInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFilestoreInstance_filestoreInstanceBasicExample(context),
			},
			{
				ResourceName:            "google_filestore_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "zone", "location"},
			},
		},
	})
}

func testAccFilestoreInstance_filestoreInstanceBasicExample(context map[string]interface{}) string {
	return Nprintf(`
resource "google_filestore_instance" "instance" {
  name = "tf-test-test-instance%{random_suffix}"
  location = "us-central1-b"
  tier = "PREMIUM"

  file_shares {
    capacity_gb = 2660
    name        = "share1"
  }

  networks {
    network = "default"
    modes   = ["MODE_IPV4"]
  }
}
`, context)
}

func TestAccFilestoreInstance_filestoreInstanceFullExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckFilestoreInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFilestoreInstance_filestoreInstanceFullExample(context),
			},
			{
				ResourceName:            "google_filestore_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "zone", "location"},
			},
		},
	})
}

func testAccFilestoreInstance_filestoreInstanceFullExample(context map[string]interface{}) string {
	return Nprintf(`
resource "google_filestore_instance" "instance" {
  name = "tf-test-test-instance%{random_suffix}"
  location = "us-central1-b"
  tier = "BASIC_SSD"

  file_shares {
    capacity_gb = 2660
    name        = "share1"

    nfs_export_options {
      ip_ranges = ["10.0.0.0/24"]
      access_mode = "READ_WRITE"
      squash_mode = "NO_ROOT_SQUASH"
   }

   nfs_export_options {
      ip_ranges = ["10.10.0.0/24"]
      access_mode = "READ_ONLY"
      squash_mode = "ROOT_SQUASH"
      anon_uid = 123
      anon_gid = 456
   }
  }

  networks {
    network = "default"
    modes   = ["MODE_IPV4"]
    connect_mode = "DIRECT_PEERING"
  }
}
`, context)
}

func testAccCheckFilestoreInstanceDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_filestore_instance" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{FilestoreBasePath}}projects/{{project}}/locations/{{location}}/instances/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(config, "GET", billingProject, url, config.UserAgent, nil, transport_tpg.IsNotFilestoreQuotaError)
			if err == nil {
				return fmt.Errorf("FilestoreInstance still exists at %s", url)
			}
		}

		return nil
	}
}
