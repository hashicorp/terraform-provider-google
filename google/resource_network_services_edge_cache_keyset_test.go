package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkServicesEdgeCacheKeyset_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkServicesEdgeCacheKeysetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkServicesEdgeCacheKeyset_networkServicesEdgeCacheKeysetBasicExample(context),
			},
			{
				ResourceName:            "google_network_services_edge_cache_keyset.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
			{
				Config: testAccNetworkServicesEdgeCacheKeyset_update(context),
			},
			{
				ResourceName:            "google_network_services_edge_cache_keyset.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
		},
	})
}

func testAccNetworkServicesEdgeCacheKeyset_update(context map[string]interface{}) string {
	return Nprintf(`

resource "google_network_services_edge_cache_keyset" "default" {
  name                 = "default%{random_suffix}"
  description          = "T2"
  public_key {
    id = "my-public-key-2"
    value = "hzd03llxB1u5FOLKFkZ6_wCJqC7jtN0bg7xlBqS6WVM"
  }
	labels = {
		a = "a"
	}
}
`, context)
}
