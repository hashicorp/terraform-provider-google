package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOrganizationIamMember_importBasic(t *testing.T) {
	t.Parallel()

	orgId := getTestOrgFromEnv(t)
	account := acctest.RandomWithPrefix("tf-test")
	projectId := getTestProjectFromEnv()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOrganizationIamMember_basic(account, orgId),
			},

			resource.TestStep{
				ResourceName:  "google_organization_iam_member.foo",
				ImportStateId: fmt.Sprintf("%s roles/browser serviceAccount:%s@%s.iam.gserviceaccount.com", orgId, account, projectId),
				ImportState:   true,
			},
		},
	})
}

func TestAccOrganizationIamBinding_importBasic(t *testing.T) {
	t.Parallel()

	orgId := getTestOrgFromEnv(t)
	account := acctest.RandomWithPrefix("tf-test")
	roleId := "tfIamTest" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOrganizationIamBinding_basic(account, roleId, orgId),
			},

			resource.TestStep{
				ResourceName:  "google_organization_iam_binding.foo",
				ImportStateId: fmt.Sprintf("%s organizations/%s/roles/%s", orgId, orgId, roleId),
				ImportState:   true,
			},
		},
	})
}
