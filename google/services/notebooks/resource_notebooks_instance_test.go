// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package notebooks_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccNotebooksInstance_create_vm_image(t *testing.T) {
	t.Parallel()

	prefix := fmt.Sprintf("%d", acctest.RandInt(t))
	name := fmt.Sprintf("tf-%s", prefix)

	acctest.VcrTest(t, resource.TestCase{
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNotebooksInstance_create_vm_image(name),
			},
			{
				ResourceName:            "google_notebooks_instance.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"vm_image", "metadata"},
			},
		},
	})
}

func TestAccNotebooksInstance_update(t *testing.T) {
	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNotebooksInstance_basic(context),
			},
			{
				ResourceName:            "google_notebooks_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"vm_image", "metadata"},
			},
			{
				Config: testAccNotebooksInstance_update(context, true),
			},
			{
				ResourceName:            "google_notebooks_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"vm_image", "metadata"},
			},
			{
				Config: testAccNotebooksInstance_update(context, false),
			},
			{
				ResourceName:            "google_notebooks_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"vm_image", "metadata"},
			},
		},
	})
}

func testAccNotebooksInstance_create_vm_image(name string) string {
	return fmt.Sprintf(`

resource "google_notebooks_instance" "test" {
  name = "%s"
  location = "us-west1-a"
  machine_type = "e2-medium"
  metadata = {
    proxy-mode = "service_account"
    terraform  = "true"
  }

  nic_type = "VIRTIO_NET"

  reservation_affinity {
    consume_reservation_type = "NO_RESERVATION"
  }

  vm_image {
    project      = "deeplearning-platform-release"
    image_family = "tf-latest-cpu"
  }
}
`, name)
}

func testAccNotebooksInstance_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_notebooks_instance" "instance" {
  name = "tf-test-notebooks-instance%{random_suffix}"
  location = "us-central1-a"
  machine_type = "e2-medium"

  vm_image {
    project      = "deeplearning-platform-release"
    image_family = "tf-latest-cpu"
  }

  metadata = {
    proxy-mode = "service_account"
    terraform  = "true"
  }

  lifecycle {
  	prevent_destroy = true
  }
}
`, context)
}

func testAccNotebooksInstance_update(context map[string]interface{}, preventDestroy bool) string {
	context["prevent_destroy"] = strconv.FormatBool(preventDestroy)

	return acctest.Nprintf(`
resource "google_notebooks_instance" "instance" {
  name = "tf-test-notebooks-instance%{random_suffix}"
  location = "us-central1-a"
  machine_type = "e2-medium"

  vm_image {
    project      = "deeplearning-platform-release"
    image_family = "tf-latest-cpu"
  }

  metadata = {
    proxy-mode = "service_account"
    terraform  = "true"
    notebook-upgrade-schedule = "0 * * * *"
  }

  labels = {
    key = "value"
  }

  lifecycle {
  	prevent_destroy = %{prevent_destroy}
  }
}
`, context)
}
