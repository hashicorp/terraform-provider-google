// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package cloudquotas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceGoogleQuotaInfo_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.google_cloud_quotas_quota_info.my_quota_info"
	service := "compute.googleapis.com"
	quotaId := "CPUS-per-project-region"

	context := map[string]interface{}{
		"project":  envvar.GetTestProjectFromEnv(),
		"service":  service,
		"quota_id": quotaId,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleQuotaInfo_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttr(resourceName, "quota_id", quotaId),
					resource.TestCheckResourceAttr(resourceName, "metric", "compute.googleapis.com/cpus"),
					resource.TestCheckResourceAttr(resourceName, "service", service),
					resource.TestCheckResourceAttrSet(resourceName, "is_precise"),
					resource.TestCheckResourceAttr(resourceName, "container_type", "PROJECT"),
					resource.TestCheckResourceAttr(resourceName, "dimensions.0", "region"),
					resource.TestCheckResourceAttr(resourceName, "metric_display_name", "CPUs"),
					resource.TestCheckResourceAttr(resourceName, "quota_display_name", "CPUs"),
					resource.TestCheckResourceAttrSet(resourceName, "metric_unit"),
					resource.TestCheckResourceAttrSet(resourceName, "quota_increase_eligibility.0.is_eligible"),
					resource.TestCheckResourceAttrSet(resourceName, "dimensions_infos.0.dimensions.region"),
					resource.TestCheckResourceAttrSet(resourceName, "dimensions_infos.0.details.0.value"),
					resource.TestCheckResourceAttrSet(resourceName, "dimensions_infos.0.applicable_locations.0"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleQuotaInfo_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
		data "google_cloud_quotas_quota_info" "my_quota_info" {
			parent      = "projects/%{project}"	
			quota_id    = "%{quota_id}"
			service 	= "%{service}"
		}
	`, context)
}
