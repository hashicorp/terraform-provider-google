package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccIAMWorkforcePoolWorkforcePool_full(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        getTestOrgFromEnv(t),
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIAMWorkforcePoolWorkforcePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIAMWorkforcePoolWorkforcePool_full(context),
			},
			{
				ResourceName:      "google_iam_workforce_pool.my_pool",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIAMWorkforcePoolWorkforcePool_update(context),
			},
			{
				ResourceName:      "google_iam_workforce_pool.my_pool",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccIAMWorkforcePoolWorkforcePool_minimal(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        getTestOrgFromEnv(t),
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIAMWorkforcePoolWorkforcePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIAMWorkforcePoolWorkforcePool_minimal(context),
			},
			{
				ResourceName:      "google_iam_workforce_pool.my_pool",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIAMWorkforcePoolWorkforcePool_update(context),
			},
			{
				ResourceName:      "google_iam_workforce_pool.my_pool",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccIAMWorkforcePoolWorkforcePool_full(context map[string]interface{}) string {
	return Nprintf(`
resource "google_iam_workforce_pool" "my_pool" {
  workforce_pool_id = "my-pool-%{random_suffix}"
  parent            = "organizations/%{org_id}"
  location          = "global"
  display_name      = "Display name"
  description       = "A sample workforce pool."
  disabled          = false
  session_duration  = "7200s"
}
`, context)
}

func testAccIAMWorkforcePoolWorkforcePool_minimal(context map[string]interface{}) string {
	return Nprintf(`
resource "google_iam_workforce_pool" "my_pool" {
  workforce_pool_id = "my-pool-%{random_suffix}"
  parent            = "organizations/%{org_id}"
  location          = "global"
}
`, context)
}

func testAccIAMWorkforcePoolWorkforcePool_update(context map[string]interface{}) string {
	return Nprintf(`
resource "google_iam_workforce_pool" "my_pool" {
  workforce_pool_id = "my-pool-%{random_suffix}"
  parent            = "organizations/%{org_id}"
  location          = "global"
  display_name      = "New display name"
  description       = "A sample workforce pool with updated description."
  disabled          = true
  session_duration  = "3600s"
}
`, context)
}
