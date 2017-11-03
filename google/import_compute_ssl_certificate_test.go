package google

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccComputeSslCertificate_import(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeSslCertificateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeSslCertificate_import,
			},
			resource.TestStep{
				ResourceName:            "google_compute_ssl_certificate.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"private_key"},
			},
		},
	})
}
