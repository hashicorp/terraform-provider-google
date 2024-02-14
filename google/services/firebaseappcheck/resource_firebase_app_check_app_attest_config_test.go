// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package firebaseappcheck_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccFirebaseAppCheckAppAttestConfig_firebaseAppCheckAppAttestConfigUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_id":    envvar.GetTestProjectFromEnv(),
		"team_id":       "9987654321",
		"random_suffix": acctest.RandString(t, 10),
		"token_ttl":     "7200s",
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
				Config: testAccFirebaseAppCheckAppAttestConfig_firebaseAppCheckAppAttestConfigMinimalExample(context),
			},
			{
				ResourceName:            "google_firebase_app_check_app_attest_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"app_id"},
			},
			{
				Config: testAccFirebaseAppCheckAppAttestConfig_firebaseAppCheckAppAttestConfigFullExample(context),
			},
			{
				ResourceName:            "google_firebase_app_check_app_attest_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"app_id"},
			},
			{
				Config: testAccFirebaseAppCheckAppAttestConfig_firebaseAppCheckAppAttestConfigMinimalExample(context),
			},
			{
				ResourceName:            "google_firebase_app_check_app_attest_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"app_id"},
			},
		},
	})
}
