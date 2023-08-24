// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package privateca_test

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/services/privateca"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func TestAccPrivatecaCaPoolIamMemberAllAuthenticatedUsersCasing(t *testing.T) {
	t.Parallel()

	capool := "tf-test-pool-iam-" + acctest.RandString(t, 10)
	project := envvar.GetTestProjectFromEnv()
	region := envvar.GetTestRegionFromEnv()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivatecaCaPoolIamMember_allAuthenticatedUsers(capool, region, project),
				Check: testAccCheckPrivatecaCaPoolIam(t, capool, region, project, "roles/privateca.certificateManager", []string{
					fmt.Sprintf("group:%s.svc.id.goog:/allAuthenticatedUsers/", project),
				}),
			},
		},
	})
}

func testAccCheckPrivatecaCaPoolIam(t *testing.T, capool, region, project, role string, members []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		d := &tpgresource.ResourceDataMock{
			FieldsInSchema: map[string]interface{}{
				"ca_pool": capool,
				"role":    role,
				"member":  "",
			},
		}
		u := &privateca.PrivatecaCaPoolIamUpdater{
			Config: acctest.GoogleProviderConfig(t),
		}
		u.SetProject(project)
		u.SetLocation(region)
		u.SetCaPool(capool)
		u.SetResourceData(d)
		p, err := u.GetResourceIamPolicy()
		if err != nil {
			return err
		}

		for _, binding := range p.Bindings {
			if binding.Role == role {
				sort.Strings(members)
				sort.Strings(binding.Members)

				if reflect.DeepEqual(members, binding.Members) {
					return nil
				}

				return fmt.Errorf("Binding found but expected members is %v, got %v", members, binding.Members)
			}
		}

		return fmt.Errorf("No binding for role %q", role)
	}
}

func testAccPrivatecaCaPoolIamMember_allAuthenticatedUsers(capool, region, project string) string {
	return fmt.Sprintf(`
resource "google_privateca_ca_pool" "default" {
  name     = "%s"
  location = "%s"
  tier     = "ENTERPRISE"
  publishing_options {
    publish_ca_cert = true
    publish_crl     = true
  }
  labels = {
    foo = "bar"
  }
}

resource "google_privateca_ca_pool_iam_member" "member" {
  ca_pool  = google_privateca_ca_pool.default.id
  role     = "roles/privateca.certificateManager"
  member   = "group:%s.svc.id.goog:/allAuthenticatedUsers/"
}
  
`, capool, region, project)
}
