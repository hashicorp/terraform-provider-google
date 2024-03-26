// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package cloudquotas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceGoogleQuotaInfos_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.google_cloud_quotas_quota_infos.my_quota_infos"
	service := "compute.googleapis.com"

	context := map[string]interface{}{
		"project": envvar.GetTestProjectFromEnv(),
		"service": service,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleQuotaInfos_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "quota_infos.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "quota_infos.0.quota_id"),
					resource.TestCheckResourceAttrSet(resourceName, "quota_infos.0.metric"),
					resource.TestCheckResourceAttr(resourceName, "quota_infos.0.service", service),
					resource.TestCheckResourceAttrSet(resourceName, "quota_infos.0.is_precise"),
					resource.TestCheckResourceAttrSet(resourceName, "quota_infos.0.container_type"),
					resource.TestCheckResourceAttrSet(resourceName, "quota_infos.0.quota_increase_eligibility.0.is_eligible"),
					resource.TestCheckResourceAttrSet(resourceName, "quota_infos.0.dimensions_infos.0.details.0.value"),
					resource.TestCheckResourceAttrSet(resourceName, "quota_infos.0.dimensions_infos.0.applicable_locations.0"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleQuotaInfos_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
		data "google_cloud_quotas_quota_infos" "my_quota_infos" {
			parent	= "projects/%{project}"	
			service	= "%{service}"
		}
	`, context)
}
