package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataLossPreventionStoredInfoType_dlpStoredInfoTypeUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       getTestProjectFromEnv(),
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataLossPreventionStoredInfoTypeDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionStoredInfoType_dlpStoredInfoTypeStart(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_stored_info_type.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataLossPreventionStoredInfoType_dlpStoredInfoTypeUpdate(context),
			},
			{
				ResourceName:      "google_data_loss_prevention_stored_info_type.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDataLossPreventionStoredInfoType_dlpStoredInfoTypeStart(context map[string]interface{}) string {
	return Nprintf(`
resource "google_data_loss_prevention_stored_info_type" "basic" {
	parent = "projects/%{project}"
	description = "Description"
	display_name = "Displayname"

	regex {
		pattern = "patient"
		group_indexes = [2]
	}
}
`, context)
}

func testAccDataLossPreventionStoredInfoType_dlpStoredInfoTypeUpdate(context map[string]interface{}) string {
	return Nprintf(`
resource "google_data_loss_prevention_stored_info_type" "basic" {
	parent = "projects/%{project}"
	description = "Updated Description"
	display_name = "display_name"

	dictionary {
		word_list {
			words = ["word", "word2"]
		}
	}
}
`, context)
}
