// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package tags_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// Tags tests cannot be run in parallel without running into Error Code 10: ABORTED
// See https://github.com/hashicorp/terraform-provider-google/issues/8637

func TestAccTags(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		"tagKeyBasic":                       testAccTagsTagKey_tagKeyBasic,
		"tagKeyBasicWithPurposeGceFirewall": testAccTagsTagKey_tagKeyBasicWithPurposeGceFirewall,
		"tagKeyUpdate":                      testAccTagsTagKey_tagKeyUpdate,
		"tagKeyIamBinding":                  testAccTagsTagKeyIamBinding,
		"tagKeyIamMember":                   testAccTagsTagKeyIamMember,
		"tagKeyIamPolicy":                   testAccTagsTagKeyIamPolicy,
		"tagValueBasic":                     testAccTagsTagValue_tagValueBasic,
		"tagValueUpdate":                    testAccTagsTagValue_tagValueUpdate,
		"tagBindingBasic":                   testAccTagsTagBinding_tagBindingBasic,
		"tagValueIamBinding":                testAccTagsTagValueIamBinding,
		"tagValueIamMember":                 testAccTagsTagValueIamMember,
		"tagValueIamPolicy":                 testAccTagsTagValueIamPolicy,
		"tagsLocationTagBindingBasic":       testAccTagsLocationTagBinding_locationTagBindingbasic,
		"tagsLocationTagBindingZonal":       TestAccTagsLocationTagBinding_locationTagBindingzonal,
	}

	for name, tc := range testCases {
		// shadow the tc variable into scope so that when
		// the loop continues, if t.Run hasn't executed tc(t)
		// yet, we don't have a race condition
		// see https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc(t)
		})
	}
}

func testAccTagsTagKey_tagKeyBasic(t *testing.T) {
	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckTagsTagKeyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTagsTagKey_tagKeyBasicExample(context),
			},
		},
	})
}

func testAccTagsTagKey_tagKeyBasicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {

  parent = "organizations/%{org_id}"
  short_name = "foo%{random_suffix}"
  description = "For foo%{random_suffix} resources."
}
`, context)
}

func testAccTagsTagKey_tagKeyBasicWithPurposeGceFirewall(t *testing.T) {
	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckTagsTagKeyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTagsTagKey_tagKeyBasicWithPurposeGceFirewallExample(context),
			},
		},
	})
}

func testAccTagsTagKey_tagKeyBasicWithPurposeGceFirewallExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "tag_network" {
	name = "vpc-%{random_suffix}"
	auto_create_subnetworks = false
}

resource "google_tags_tag_key" "key" {
	  parent = "organizations/%{org_id}"
	  short_name = "foo%{random_suffix}"
	  description = "For foo%{random_suffix} resources."
	  purpose = "GCE_FIREWALL"
	  # purpose_data expects either a selfLinkWithId (not a property of google_compute_network) or the format <project-name>/<vpc-name>.
	  # selfLink is not sufficient and will result in an error, so we build a string to match the second option.
	  purpose_data = {network = "${google_compute_network.tag_network.project}/${google_compute_network.tag_network.name}"}
	}

`, context)
}

func testAccTagsTagKey_tagKeyUpdate(t *testing.T) {
	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckTagsTagKeyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTagsTagKey_basic(context),
			},
			{
				ResourceName:      "google_tags_tag_key.key",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTagsTagKey_basicUpdated(context),
			},
			{
				ResourceName:      "google_tags_tag_key.key",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccTagsTagKey_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {

  parent = "organizations/%{org_id}"
  short_name = "foo%{random_suffix}"
  description = "For foo%{random_suffix} resources."
}
`, context)
}

func testAccTagsTagKey_basicUpdated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {

  parent = "organizations/%{org_id}"
  short_name = "foo%{random_suffix}"
  description = "Anything related to foo%{random_suffix}"
}
`, context)
}

func testAccCheckTagsTagKeyDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_tags_tag_key" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{TagsBasePath}}tagKeys/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("TagsTagKey still exists at %s", url)
			}
		}

		return nil
	}
}

func testAccTagsTagValue_tagValueBasic(t *testing.T) {
	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckTagsTagValueDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTagsTagValue_tagValueBasicExample(context),
			},
		},
	})
}

func testAccTagsTagValue_tagValueBasicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {

  parent = "organizations/%{org_id}"
  short_name = "foobarbaz%{random_suffix}"
  description = "For foo/bar/baz resources."
}

resource "google_tags_tag_value" "value" {

  parent = "tagKeys/${google_tags_tag_key.key.name}"
  short_name = "foo%{random_suffix}"
  description = "For foo resources."
}
`, context)
}

func testAccTagsTagValue_tagValueUpdate(t *testing.T) {
	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckTagsTagValueDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTagsTagValue_basic(context),
			},
			{
				ResourceName:      "google_tags_tag_key.key",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTagsTagValue_basicUpdated(context),
			},
			{
				ResourceName:      "google_tags_tag_key.key",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccTagsTagValue_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {

  parent = "organizations/%{org_id}"
  short_name = "foobarbaz%{random_suffix}"
  description = "For foo/bar/baz resources."
}

resource "google_tags_tag_value" "value" {

  parent = "tagKeys/${google_tags_tag_key.key.name}"
  short_name = "foo%{random_suffix}"
  description = "For foo resources."
}
`, context)
}

func testAccTagsTagValue_basicUpdated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {

  parent = "organizations/%{org_id}"
  short_name = "foobarbaz%{random_suffix}"
  description = "For foo/bar/baz resources."
}

resource "google_tags_tag_value" "value" {

  parent = "tagKeys/${google_tags_tag_key.key.name}"
  short_name = "foo%{random_suffix}"
  description = "For any foo resources."
}
`, context)
}

func testAccCheckTagsTagValueDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_tags_tag_key" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{TagsBasePath}}tagValues/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("TagsTagValue still exists at %s", url)
			}
		}

		return nil
	}
}

func testAccTagsTagBinding_tagBindingBasic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"project_id":    "tf-test-" + acctest.RandString(t, 10),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
		},
		CheckDestroy: testAccCheckTagsTagBindingDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTagsTagBinding_tagBindingBasicExample(context),
			},
		},
	})
}

func testAccTagsTagBinding_tagBindingBasicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
	project_id = "%{project_id}"
	name       = "%{project_id}"
	org_id     = "%{org_id}"
}

resource "google_tags_tag_key" "key" {
	parent = "organizations/%{org_id}"
	short_name = "keyname%{random_suffix}"
	description = "For a certain set of resources."
}

resource "google_tags_tag_value" "value" {
	parent = "tagKeys/${google_tags_tag_key.key.name}"
	short_name = "foo%{random_suffix}"
	description = "For foo%{random_suffix} resources."
}

resource "google_tags_tag_binding" "binding" {
	parent = "//cloudresourcemanager.googleapis.com/projects/${google_project.project.number}"
	tag_value = "tagValues/${google_tags_tag_value.value.name}"
}
`, context)
}

func testAccCheckTagsTagBindingDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_tags_tag_binding" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{TagsBasePath}}tagBindings/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("TagsTagBinding still exists at %s", url)
			}
		}

		return nil
	}
}

func testAccTagsTagKeyIamBinding(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"role":          "roles/viewer",
		"org_id":        envvar.GetTestOrgFromEnv(t),

		"short_name": "tf-test-key-" + acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTagsTagKeyIamBinding_basicGenerated(context),
			},
			{
				// Test Iam Binding update
				Config: testAccTagsTagKeyIamBinding_updateGenerated(context),
			},
		},
	})
}

func testAccTagsTagKeyIamMember(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"role":          "roles/viewer",
		"org_id":        envvar.GetTestOrgFromEnv(t),

		"short_name": "tf-test-key-" + acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test Iam Member creation (no update for member, no need to test)
				Config: testAccTagsTagKeyIamMember_basicGenerated(context),
			},
		},
	})
}

func testAccTagsTagKeyIamPolicy(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"role":          "roles/viewer",
		"org_id":        envvar.GetTestOrgFromEnv(t),

		"short_name": "tf-test-key-" + acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTagsTagKeyIamPolicy_basicGenerated(context),
			},
			{
				Config: testAccTagsTagKeyIamPolicy_emptyBinding(context),
			},
		},
	})
}

func testAccTagsTagKeyIamMember_basicGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {

  parent = "organizations/%{org_id}"
  short_name = "%{short_name}"
  description = "For %{short_name} resources."
}

resource "google_tags_tag_key_iam_member" "foo" {
  tag_key = google_tags_tag_key.key.name
  role = "%{role}"
  member = "user:admin@hashicorptest.com"
}
`, context)
}

func testAccTagsTagKeyIamPolicy_basicGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {

  parent = "organizations/%{org_id}"
  short_name = "%{short_name}"
  description = "For %{short_name} resources."
}

data "google_iam_policy" "foo" {
  binding {
    role = "%{role}"
    members = ["user:admin@hashicorptest.com"]
  }
}

resource "google_tags_tag_key_iam_policy" "foo" {
  tag_key = google_tags_tag_key.key.name
  policy_data = data.google_iam_policy.foo.policy_data
}
`, context)
}

func testAccTagsTagKeyIamPolicy_emptyBinding(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {

  parent = "organizations/%{org_id}"
  short_name = "%{short_name}"
  description = "For %{short_name} resources."
}

data "google_iam_policy" "foo" {
}

resource "google_tags_tag_key_iam_policy" "foo" {
  tag_key = google_tags_tag_key.key.name
  policy_data = data.google_iam_policy.foo.policy_data
}
`, context)
}

func testAccTagsTagKeyIamBinding_basicGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {

  parent = "organizations/%{org_id}"
  short_name = "%{short_name}"
  description = "For %{short_name} resources."
}

resource "google_tags_tag_key_iam_binding" "foo" {
  tag_key = google_tags_tag_key.key.name
  role = "%{role}"
  members = ["user:admin@hashicorptest.com"]
}
`, context)
}

func testAccTagsTagKeyIamBinding_updateGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {

  parent = "organizations/%{org_id}"
  short_name = "%{short_name}"
  description = "For %{short_name} resources."
}

resource "google_tags_tag_key_iam_binding" "foo" {
  tag_key = google_tags_tag_key.key.name
  role = "%{role}"
  members = ["user:admin@hashicorptest.com", "user:gterraformtest1@gmail.com"]
}
`, context)
}

func testAccTagsTagValueIamBinding(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"role":          "roles/viewer",
		"org_id":        envvar.GetTestOrgFromEnv(t),

		"key_short_name":   "tf-test-key-" + acctest.RandString(t, 10),
		"value_short_name": "tf-test-value-" + acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTagsTagValueIamBinding_basicGenerated(context),
			},
			{
				// Test Iam Binding update
				Config: testAccTagsTagValueIamBinding_updateGenerated(context),
			},
		},
	})
}

func testAccTagsTagValueIamMember(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"role":          "roles/viewer",
		"org_id":        envvar.GetTestOrgFromEnv(t),

		"key_short_name":   "tf-test-key-" + acctest.RandString(t, 10),
		"value_short_name": "tf-test-value-" + acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test Iam Member creation (no update for member, no need to test)
				Config: testAccTagsTagValueIamMember_basicGenerated(context),
			},
		},
	})
}

func testAccTagsTagValueIamPolicy(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"role":          "roles/viewer",
		"org_id":        envvar.GetTestOrgFromEnv(t),

		"key_short_name":   "tf-test-key-" + acctest.RandString(t, 10),
		"value_short_name": "tf-test-value-" + acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTagsTagValueIamPolicy_basicGenerated(context),
			},
			{
				Config: testAccTagsTagValueIamPolicy_emptyBinding(context),
			},
		},
	})
}

func testAccTagsTagValueIamMember_basicGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {
	parent = "organizations/%{org_id}"
	short_name = "%{key_short_name}"
	description = "For %{key_short_name} resources."
}

resource "google_tags_tag_value" "value" {
	parent = "tagKeys/${google_tags_tag_key.key.name}"
	short_name = "%{value_short_name}"
	description = "For %{value_short_name} resources."
}

resource "google_tags_tag_value_iam_member" "foo" {
  tag_value = google_tags_tag_value.value.name
  role = "%{role}"
  member = "user:admin@hashicorptest.com"
}
`, context)
}

func testAccTagsTagValueIamPolicy_basicGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {
	parent = "organizations/%{org_id}"
	short_name = "%{key_short_name}"
	description = "For %{key_short_name} resources."
}

resource "google_tags_tag_value" "value" {
	parent = "tagKeys/${google_tags_tag_key.key.name}"
	short_name = "%{value_short_name}"
	description = "For %{value_short_name} resources."
}

data "google_iam_policy" "foo" {
  binding {
    role = "%{role}"
    members = ["user:admin@hashicorptest.com"]
  }
}

resource "google_tags_tag_value_iam_policy" "foo" {
  tag_value = google_tags_tag_value.value.name
  policy_data = data.google_iam_policy.foo.policy_data
}
`, context)
}

func testAccTagsTagValueIamPolicy_emptyBinding(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {
	parent = "organizations/%{org_id}"
	short_name = "%{key_short_name}"
	description = "For %{key_short_name} resources."
}

resource "google_tags_tag_value" "value" {
	parent = "tagKeys/${google_tags_tag_key.key.name}"
	short_name = "%{value_short_name}"
	description = "For %{value_short_name} resources."
}

data "google_iam_policy" "foo" {
}

resource "google_tags_tag_value_iam_policy" "foo" {
  tag_value = google_tags_tag_value.value.name
  policy_data = data.google_iam_policy.foo.policy_data
}
`, context)
}

func testAccTagsTagValueIamBinding_basicGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {
	parent = "organizations/%{org_id}"
	short_name = "%{key_short_name}"
	description = "For %{key_short_name} resources."
}

resource "google_tags_tag_value" "value" {
	parent = "tagKeys/${google_tags_tag_key.key.name}"
	short_name = "%{value_short_name}"
	description = "For %{value_short_name} resources."
}

resource "google_tags_tag_value_iam_binding" "foo" {
  tag_value = google_tags_tag_value.value.name
  role = "%{role}"
  members = ["user:admin@hashicorptest.com"]
}
`, context)
}

func testAccTagsTagValueIamBinding_updateGenerated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key" {
	parent = "organizations/%{org_id}"
	short_name = "%{key_short_name}"
	description = "For %{key_short_name} resources."
}

resource "google_tags_tag_value" "value" {
	parent = "tagKeys/${google_tags_tag_key.key.name}"
	short_name = "%{value_short_name}"
	description = "For %{value_short_name} resources."
}

resource "google_tags_tag_value_iam_binding" "foo" {
  tag_value = google_tags_tag_value.value.name
  role = "%{role}"
  members = ["user:admin@hashicorptest.com", "user:gterraformtest1@gmail.com"]
}
`, context)
}

func testAccTagsLocationTagBinding_locationTagBindingbasic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		// "org_id":        envvar.GetTestOrgFromEnv(t),
		// "project_id":    "tf-test-" + acctest.RandString(t, 10),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
		},
		CheckDestroy: testAccCheckTagsLocationTagBindingDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTagsLocationTagBinding_locationTagBindingBasicExample(context),
			},
			{
				ResourceName:      "google_tags_location_tag_binding.binding",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccTagsLocationTagBinding_locationTagBindingBasicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

resource "google_tags_tag_key" "key" {
	parent = "organizations/${data.google_project.project.org_id}"
	short_name = "keyname%{random_suffix}"
	description = "For a certain set of resources."
}

resource "google_tags_tag_value" "value" {
	parent = "tagKeys/${google_tags_tag_key.key.name}"
	short_name = "foo%{random_suffix}"
	description = "For foo%{random_suffix} resources."
}

resource "google_cloud_run_service" "default" {
	name     = "tf-test-cloudrun-srv%{random_suffix}"
	location = "us-central1"
  
	template {
	  spec {
		containers {
		  image = "us-docker.pkg.dev/cloudrun/container/hello"
		}
	  }
	}
  
	traffic {
	  percent         = 100
	  latest_revision = true
	}
}
  
resource "google_tags_location_tag_binding" "binding" {
	parent = "//run.googleapis.com/projects/${data.google_project.project.number}/locations/${google_cloud_run_service.default.location}/services/${google_cloud_run_service.default.name}"
	tag_value = "tagValues/${google_tags_tag_value.value.name}"
	location = "us-central1"
}
`, context)
}

func testAccCheckTagsLocationTagBindingDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_tags_location_tag_binding" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{TagsLocationBasePath}}{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("TagsTagBinding still exists at %s", url)
			}
		}
		return nil
	}
}

func TestAccTagsLocationTagBinding_locationTagBindingzonal(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
		},
		CheckDestroy: testAccCheckTagsLocationTagBindingDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTagsLocationTagBinding_locationTagBindingZonalExample(context),
			},
			{
				ResourceName:      "google_tags_location_tag_binding.binding",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccTagsLocationTagBinding_locationTagBindingZonalExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}
resource "google_tags_tag_key" "key" {
	parent = "organizations/%{org_id}"
	short_name = "keyname%{random_suffix}"
	description = "For a certain set of resources."
}
resource "google_tags_tag_value" "value" {
	parent = "tagKeys/${google_tags_tag_key.key.name}"
	short_name = "foo%{random_suffix}"
	description = "For foo%{random_suffix} resources."
}
resource "google_compute_instance" "default" {
	name         = "test-%{random_suffix}"
	machine_type = "e2-medium"
	zone         = "us-central1-a"
	boot_disk {
		initialize_params {
			image = "debian-cloud/debian-11"
		}
	}
	network_interface {
		 network = "default"
	}
}
resource "google_tags_location_tag_binding" "binding" {
	parent = "//compute.googleapis.com/projects/${data.google_project.project.number}/zones/us-central1-a/instances/${google_compute_instance.default.instance_id}"
	tag_value = "tagValues/${google_tags_tag_value.value.name}"
	location = "us-central1-a"
}
`, context)
}
