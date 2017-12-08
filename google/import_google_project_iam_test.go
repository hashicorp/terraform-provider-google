package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccGoogleProjectIamMember_importBasic(t *testing.T) {
	t.Parallel()

	resourceName := "google_project_iam_member.acceptance"
	org := getTestOrgFromEnv(t)
	pid := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccGoogleProjectAssociateMemberBasic(pid, "Acceptance", org),
			},

			resource.TestStep{
				ResourceName:  resourceName,
				ImportStateId: fmt.Sprintf("%s %s %s", pid, "roles/compute.instanceAdmin", "user:admin@hashicorptest.com"),
				ImportState:   true,
			},
		},
	})
}

func TestAccGoogleProjectIamBinding_importBasic(t *testing.T) {
	t.Parallel()

	resourceName := "google_project_iam_binding.acceptance"
	org := getTestOrgFromEnv(t)
	pid := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccGoogleProjectAssociateBindingBasic(pid, "Acceptance", org),
			},

			resource.TestStep{
				ResourceName:  resourceName,
				ImportStateId: fmt.Sprintf("%s %s", pid, "roles/compute.instanceAdmin"),
				ImportState:   true,
			},
		},
	})
}
