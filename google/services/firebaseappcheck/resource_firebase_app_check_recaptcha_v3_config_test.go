// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package firebaseappcheck_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccFirebaseAppCheckRecaptchaV3Config_firebaseAppCheckRecaptchaV3ConfigUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_id":    envvar.GetTestProjectFromEnv(),
		"token_ttl":     "7200s",
		"site_secret":   "6Lf9YnQpAAAAAC3-MHmdAllTbPwTZxpUw5d34YzX",
		"random_suffix": acctest.RandString(t, 10),
	}

	contextUpdated := map[string]interface{}{
		"project_id":    envvar.GetTestProjectFromEnv(),
		"token_ttl":     "3800s",
		"site_secret":   "7Lf9YnQpAAAAAC3-MHmdAllTbPwTZxpUw5d34YzX",
		"random_suffix": context["random_suffix"],
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderBetaFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
			"time":   {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccFirebaseAppCheckRecaptchaV3Config_firebaseAppCheckRecaptchaV3ConfigBasicExample(context),
			},
			{
				ResourceName:            "google_firebase_app_check_recaptcha_v3_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"site_secret", "app_id"},
			},
			{
				Config: testAccFirebaseAppCheckRecaptchaV3Config_firebaseAppCheckRecaptchaV3ConfigBasicExample(contextUpdated),
			},
			{
				ResourceName:            "google_firebase_app_check_recaptcha_v3_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"site_secret", "app_id"},
			},
		},
	})
}
