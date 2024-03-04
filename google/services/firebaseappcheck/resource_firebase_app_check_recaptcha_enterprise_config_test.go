// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package firebaseappcheck_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccFirebaseAppCheckRecaptchaEnterpriseConfig_firebaseAppCheckRecaptchaEnterpriseConfigUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_id":    envvar.GetTestProjectFromEnv(),
		"token_ttl":     "7200s",
		"site_key":      "6LdpMXIpAAAAANkwWQPgEdjEhal7ugkH9RK9ytuw",
		"random_suffix": acctest.RandString(t, 10),
	}

	contextUpdated := map[string]interface{}{
		"project_id":    envvar.GetTestProjectFromEnv(),
		"token_ttl":     "3800s",
		"site_key":      "7LdpMXIpAAAAANkwWQPgEdjEhal7ugkH9RK9ytuw",
		"random_suffix": context["random_suffix"],
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
			"time":   {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccFirebaseAppCheckRecaptchaEnterpriseConfig_firebaseAppCheckRecaptchaEnterpriseConfigBasicExample(context),
			},
			{
				ResourceName:            "google_firebase_app_check_recaptcha_enterprise_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"app_id"},
			},
			{
				Config: testAccFirebaseAppCheckRecaptchaEnterpriseConfig_firebaseAppCheckRecaptchaEnterpriseConfigBasicExample(contextUpdated),
			},
			{
				ResourceName:            "google_firebase_app_check_recaptcha_enterprise_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"app_id"},
			},
		},
	})
}
