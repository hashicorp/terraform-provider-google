// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGoogleDatastreamStaticIps_basic(t *testing.T) {
	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleDatastreamStaticIps_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.google_datastream_static_ips.foobar",
						"static_ips.0", regexp.MustCompile("^\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}$")),
					resource.TestMatchResourceAttr("data.google_datastream_static_ips.foobarbaz",
						"static_ips.0", regexp.MustCompile("^\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}$")),
				),
			},
		},
	})
}

const testAccDataSourceGoogleDatastreamStaticIps_basic = `
data "google_project" "project" {
}

data "google_datastream_static_ips" "foobar" {
	location       = "us-west1"
}

data "google_datastream_static_ips" "foobarbaz" {
	location       = "us-central1"
	project        = data.google_project.project.project_id
}
`
