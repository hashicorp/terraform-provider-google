package google

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccComputeSharedVpc_importBasic(t *testing.T) {
	org := getTestOrgFromEnv(t)
	billingId := getTestBillingAccountFromEnv(t)

	hostProject := "xpn-host-" + acctest.RandString(10)
	serviceProject := "xpn-service-" + acctest.RandString(10)

	hostProjectResourceName := "google_compute_shared_vpc_host_project.host"
	serviceProjectResourceName := "google_compute_shared_vpc_service_project.service"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeSharedVpc_basic(hostProject, serviceProject, org, billingId),
			},

			resource.TestStep{
				ResourceName:      hostProjectResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},

			resource.TestStep{
				ResourceName:      serviceProjectResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
