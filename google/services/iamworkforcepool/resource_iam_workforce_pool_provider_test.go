// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package iamworkforcepool_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccIAMWorkforcePoolWorkforcePoolProvider_oidc(t *testing.T) {
	t.Parallel()

	random_suffix := acctest.RandString(t, 10)
	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": random_suffix,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIAMWorkforcePoolWorkforcePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIAMWorkforcePoolWorkforcePoolProvider_oidc_full(context),
			},
			{
				ResourceName:            "google_iam_workforce_pool_provider.my_provider",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"oidc.0.client_secret.0.value.0.plain_text"},
			},
			{
				Config: testAccIAMWorkforcePoolWorkforcePoolProvider_oidc_update(context),
			},
			{
				ResourceName:            "google_iam_workforce_pool_provider.my_provider",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"oidc.0.client_secret.0.value.0.plain_text"},
			},
			{
				Config: testAccIAMWorkforcePoolWorkforcePoolProvider_oidc_update_clearClientSecret(context),
			},
			{
				ResourceName:      "google_iam_workforce_pool_provider.my_provider",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIAMWorkforcePoolWorkforcePoolProvider_oidc_basic(context),
			},
			{
				ResourceName:            "google_iam_workforce_pool_provider.my_provider",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"oidc.0.client_secret.0.value.0.plain_text"},
			},
			{
				Config: testAccIAMWorkforcePoolWorkforcePoolProvider_destroy(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIAMWorkforcePoolWorkforcePoolProviderAccess(t, random_suffix),
				),
			},
		},
	})
}

func TestAccIAMWorkforcePoolWorkforcePoolProvider_saml(t *testing.T) {
	t.Parallel()

	random_suffix := acctest.RandString(t, 10)
	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": random_suffix,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIAMWorkforcePoolWorkforcePoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIAMWorkforcePoolWorkforcePoolProvider_saml_full(context),
			},
			{
				ResourceName:      "google_iam_workforce_pool_provider.my_provider",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIAMWorkforcePoolWorkforcePoolProvider_saml_update(context),
			},
			{
				ResourceName:      "google_iam_workforce_pool_provider.my_provider",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIAMWorkforcePoolWorkforcePoolProvider_saml_basic(context),
			},
			{
				ResourceName:      "google_iam_workforce_pool_provider.my_provider",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIAMWorkforcePoolWorkforcePoolProvider_destroy(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIAMWorkforcePoolWorkforcePoolProviderAccess(t, random_suffix),
				),
			},
		},
	})
}

func testAccCheckIAMWorkforcePoolWorkforcePoolProviderAccess(t *testing.T, random_suffix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		pool_resource_name := "google_iam_workforce_pool.my_pool"
		pool_rs, ok := s.RootModule().Resources[pool_resource_name]
		if !ok {
			return fmt.Errorf("Resource %s Not found", pool_resource_name)
		}
		config := acctest.GoogleProviderConfig(t)

		pool_url, err := tpgresource.ReplaceVarsForTest(config, pool_rs, "{{IAMWorkforcePoolBasePath}}locations/{{location}}/workforcePools/{{workforce_pool_id}}")
		if err != nil {
			return err
		}

		url := fmt.Sprintf("%s/providers/my-provider-%s", pool_url, random_suffix)
		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			RawURL:    url,
			UserAgent: config.UserAgent,
		})
		if err != nil {
			return nil
		}

		if v := res["state"]; v == "DELETED" {
			return nil
		}

		return fmt.Errorf("IAMWorkforcePoolProvider still exists at %s", url)
	}
}

func testAccIAMWorkforcePoolWorkforcePoolProvider_oidc_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_workforce_pool" "my_pool" {
  workforce_pool_id = "my-pool-%{random_suffix}"
  parent            = "organizations/%{org_id}"
  location          = "global"
}

resource "google_iam_workforce_pool_provider" "my_provider" {
  workforce_pool_id   = google_iam_workforce_pool.my_pool.workforce_pool_id
  location            = google_iam_workforce_pool.my_pool.location
  provider_id         = "my-provider-%{random_suffix}"
  attribute_mapping   = {
    "google.subject"  = "assertion.sub"
  }
  oidc {
    issuer_uri        = "https://accounts.thirdparty.com"
    client_id         = "client-id"
    client_secret {
      value {
        plain_text = "client-secret"
      }
    }
    web_sso_config {
      response_type             = "CODE"
      assertion_claims_behavior = "MERGE_USER_INFO_OVER_ID_TOKEN_CLAIMS"
    }
  }
  display_name        = "Display name"
  description         = "A sample OIDC workforce pool provider."
  disabled            = false
  attribute_condition = "true"
}
`, context)
}

func testAccIAMWorkforcePoolWorkforcePoolProvider_oidc_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_workforce_pool" "my_pool" {
  workforce_pool_id = "my-pool-%{random_suffix}"
  parent            = "organizations/%{org_id}"
  location          = "global"
}

resource "google_iam_workforce_pool_provider" "my_provider" {
  workforce_pool_id   = google_iam_workforce_pool.my_pool.workforce_pool_id
  location            = google_iam_workforce_pool.my_pool.location
  provider_id         = "my-provider-%{random_suffix}"
  attribute_mapping   = {
    "google.subject"  = "false"
  }
  oidc {
    issuer_uri        = "https://test.thirdparty.com"
    client_id         = "new-client-id"
    client_secret {
      value {
        plain_text = "new-client-secret"
      }
    }
    web_sso_config {
      response_type             = "ID_TOKEN"
      assertion_claims_behavior = "ONLY_ID_TOKEN_CLAIMS"
    }
  }
  display_name        = "New Display name"
  description         = "A sample OIDC workforce pool provider with updated description."
  disabled            = true
  attribute_condition = "false"
}
`, context)
}

func testAccIAMWorkforcePoolWorkforcePoolProvider_oidc_update_clearClientSecret(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_workforce_pool" "my_pool" {
  workforce_pool_id = "my-pool-%{random_suffix}"
  parent            = "organizations/%{org_id}"
  location          = "global"
}

resource "google_iam_workforce_pool_provider" "my_provider" {
  workforce_pool_id   = google_iam_workforce_pool.my_pool.workforce_pool_id
  location            = google_iam_workforce_pool.my_pool.location
  provider_id         = "my-provider-%{random_suffix}"
  attribute_mapping   = {
    "google.subject"  = "false"
  }
  oidc {
    issuer_uri        = "https://test.thirdparty.com"
    client_id         = "new-client-id"
    web_sso_config {
      response_type             = "ID_TOKEN"
      assertion_claims_behavior = "ONLY_ID_TOKEN_CLAIMS"
    }
  }
  display_name        = "New Display name"
  description         = "A sample OIDC workforce pool provider with updated description."
  disabled            = true
  attribute_condition = "false"
}
`, context)
}

func testAccIAMWorkforcePoolWorkforcePoolProvider_oidc_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_workforce_pool" "my_pool" {
  workforce_pool_id = "my-pool-%{random_suffix}"
  parent            = "organizations/%{org_id}"
  location          = "global"
}

resource "google_iam_workforce_pool_provider" "my_provider" {
  workforce_pool_id  = google_iam_workforce_pool.my_pool.workforce_pool_id
  location           = google_iam_workforce_pool.my_pool.location
  provider_id        = "my-provider-%{random_suffix}"
  attribute_mapping  = {
    "google.subject" = "assertion.sub"
  }
  oidc {
    issuer_uri       = "https://accounts.thirdparty.com"
    client_id        = "client-id"
    client_secret {
      value {
        plain_text = "client-secret"
      }
    }
    web_sso_config {
      response_type             = "CODE"
      assertion_claims_behavior = "MERGE_USER_INFO_OVER_ID_TOKEN_CLAIMS"
    }
  }
}
`, context)
}

func testAccIAMWorkforcePoolWorkforcePoolProvider_saml_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_workforce_pool" "my_pool" {
  workforce_pool_id = "my-pool-%{random_suffix}"
  parent            = "organizations/%{org_id}"
  location          = "global"
}

resource "google_iam_workforce_pool_provider" "my_provider" {
  workforce_pool_id   = google_iam_workforce_pool.my_pool.workforce_pool_id
  location            = google_iam_workforce_pool.my_pool.location
  provider_id         = "my-provider-%{random_suffix}"
  attribute_mapping   = {
    "google.subject"  = "assertion.sub"
  }
  saml {
    idp_metadata_xml  = "<?xml version=\"1.0\"?><md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\" entityID=\"https://test.com\"><md:IDPSSODescriptor protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\"> <md:KeyDescriptor use=\"signing\"><ds:KeyInfo xmlns:ds=\"http://www.w3.org/2000/09/xmldsig#\"><ds:X509Data><ds:X509Certificate>MIIDpDCCAoygAwIBAgIGAX7/5qPhMA0GCSqGSIb3DQEBCwUAMIGSMQswCQYDVQQGEwJVUzETMBEGA1UECAwKQ2FsaWZvcm5pYTEWMBQGA1UEBwwNU2FuIEZyYW5jaXNjbzENMAsGA1UECgwET2t0YTEUMBIGA1UECwwLU1NPUHJvdmlkZXIxEzARBgNVBAMMCmRldi00NTg0MjExHDAaBgkqhkiG9w0BCQEWDWluZm9Ab2t0YS5jb20wHhcNMjIwMjE2MDAxOTEyWhcNMzIwMjE2MDAyMDEyWjCBkjELMAkGA1UEBhMCVVMxEzARBgNVBAgMCkNhbGlmb3JuaWExFjAUBgNVBAcMDVNhbiBGcmFuY2lzY28xDTALBgNVBAoMBE9rdGExFDASBgNVBAsMC1NTT1Byb3ZpZGVyMRMwEQYDVQQDDApkZXYtNDU4NDIxMRwwGgYJKoZIhvcNAQkBFg1pbmZvQG9rdGEuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAxrBl7GKz52cRpxF9xCsirnRuMxnhFBaUrsHqAQrLqWmdlpNYZTVg+T9iQ+aq/iE68L+BRZcZniKIvW58wqqS0ltXVvIkXuDSvnvnkkI5yMIVErR20K8jSOKQm1FmK+fgAJ4koshFiu9oLiqu0Ejc0DuL3/XRsb4RuxjktKTb1khgBBtb+7idEk0sFR0RPefAweXImJkDHDm7SxjDwGJUubbqpdTxasPr0W+AHI1VUzsUsTiHAoyb0XDkYqHfDzhj/ZdIEl4zHQ3bEZvlD984ztAnmX2SuFLLKfXeAAGHei8MMixJvwxYkkPeYZ/5h8WgBZPP4heS2CPjwYExt29L8QIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQARjJFz++a9Z5IQGFzsZMrX2EDR5ML4xxUiQkbhld1S1PljOLcYFARDmUC2YYHOueU4ee8Jid9nPGEUebV/4Jok+b+oQh+dWMgiWjSLI7h5q4OYZ3VJtdlVwgMFt2iz+/4yBKMUZ50g3Qgg36vE34us+eKitg759JgCNsibxn0qtJgSPm0sgP2L6yTaLnoEUbXBRxCwynTSkp9ZijZqEzbhN0e2dWv7Rx/nfpohpDP6vEiFImKFHpDSv3M/5de1ytQzPFrZBYt9WlzlYwE1aD9FHCxdd+rWgYMVVoRaRmndpV/Rq3QUuDuFJtaoX11bC7ExkOpg9KstZzA63i3VcfYv</ds:X509Certificate></ds:X509Data></ds:KeyInfo></md:KeyDescriptor><md:SingleSignOnService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect\" Location=\"https://test.com/sso\"/></md:IDPSSODescriptor></md:EntityDescriptor>"
  }
  display_name        = "Display name"
  description         = "A sample SAML workforce pool provider."
  disabled            = false
  attribute_condition = "true"
}
`, context)
}

func testAccIAMWorkforcePoolWorkforcePoolProvider_saml_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_workforce_pool" "my_pool" {
  workforce_pool_id = "my-pool-%{random_suffix}"
  parent            = "organizations/%{org_id}"
  location          = "global"
}

resource "google_iam_workforce_pool_provider" "my_provider" {
  workforce_pool_id   = google_iam_workforce_pool.my_pool.workforce_pool_id
  location            = google_iam_workforce_pool.my_pool.location
  provider_id         = "my-provider-%{random_suffix}"
  attribute_mapping   = {
    "google.subject": "false"
  }
  saml {
    idp_metadata_xml  = "<?xml version=\"1.0\"?><md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\" entityID=\"https://new-test.com\"><md:IDPSSODescriptor protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\"> <md:KeyDescriptor use=\"signing\"><ds:KeyInfo xmlns:ds=\"http://www.w3.org/2000/09/xmldsig#\"><ds:X509Data><ds:X509Certificate>MIIDpDCCAoygAwIBAgIGAX7/5qPhMA0GCSqGSIb3DQEBCwUAMIGSMQswCQYDVQQGEwJVUzETMBEGA1UECAwKQ2FsaWZvcm5pYTEWMBQGA1UEBwwNU2FuIEZyYW5jaXNjbzENMAsGA1UECgwET2t0YTEUMBIGA1UECwwLU1NPUHJvdmlkZXIxEzARBgNVBAMMCmRldi00NTg0MjExHDAaBgkqhkiG9w0BCQEWDWluZm9Ab2t0YS5jb20wHhcNMjIwMjE2MDAxOTEyWhcNMzIwMjE2MDAyMDEyWjCBkjELMAkGA1UEBhMCVVMxEzARBgNVBAgMCkNhbGlmb3JuaWExFjAUBgNVBAcMDVNhbiBGcmFuY2lzY28xDTALBgNVBAoMBE9rdGExFDASBgNVBAsMC1NTT1Byb3ZpZGVyMRMwEQYDVQQDDApkZXYtNDU4NDIxMRwwGgYJKoZIhvcNAQkBFg1pbmZvQG9rdGEuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAxrBl7GKz52cRpxF9xCsirnRuMxnhFBaUrsHqAQrLqWmdlpNYZTVg+T9iQ+aq/iE68L+BRZcZniKIvW58wqqS0ltXVvIkXuDSvnvnkkI5yMIVErR20K8jSOKQm1FmK+fgAJ4koshFiu9oLiqu0Ejc0DuL3/XRsb4RuxjktKTb1khgBBtb+7idEk0sFR0RPefAweXImJkDHDm7SxjDwGJUubbqpdTxasPr0W+AHI1VUzsUsTiHAoyb0XDkYqHfDzhj/ZdIEl4zHQ3bEZvlD984ztAnmX2SuFLLKfXeAAGHei8MMixJvwxYkkPeYZ/5h8WgBZPP4heS2CPjwYExt29L8QIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQARjJFz++a9Z5IQGFzsZMrX2EDR5ML4xxUiQkbhld1S1PljOLcYFARDmUC2YYHOueU4ee8Jid9nPGEUebV/4Jok+b+oQh+dWMgiWjSLI7h5q4OYZ3VJtdlVwgMFt2iz+/4yBKMUZ50g3Qgg36vE34us+eKitg759JgCNsibxn0qtJgSPm0sgP2L6yTaLnoEUbXBRxCwynTSkp9ZijZqEzbhN0e2dWv7Rx/nfpohpDP6vEiFImKFHpDSv3M/5de1ytQzPFrZBYt9WlzlYwE1aD9FHCxdd+rWgYMVVoRaRmndpV/Rq3QUuDuFJtaoX11bC7ExkOpg9KstZzA63i3VcfYv</ds:X509Certificate></ds:X509Data></ds:KeyInfo></md:KeyDescriptor><md:SingleSignOnService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect\" Location=\"https://test.com/sso\"/></md:IDPSSODescriptor></md:EntityDescriptor>"
  }
  display_name        = "New Display name"
  description         = "A sample SAML workforce pool provider with updated description."
  disabled            = true
  attribute_condition = "false"
}
`, context)
}

func testAccIAMWorkforcePoolWorkforcePoolProvider_saml_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_workforce_pool" "my_pool" {
  workforce_pool_id = "my-pool-%{random_suffix}"
  parent            = "organizations/%{org_id}"
  location          = "global"
}

resource "google_iam_workforce_pool_provider" "my_provider" {
  workforce_pool_id  = google_iam_workforce_pool.my_pool.workforce_pool_id
  location           = google_iam_workforce_pool.my_pool.location
  provider_id        = "my-provider-%{random_suffix}"
  attribute_mapping  = {
    "google.subject" = "assertion.sub"
  }
  saml {
    idp_metadata_xml = "<?xml version=\"1.0\"?><md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\" entityID=\"https://test.com\"><md:IDPSSODescriptor protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\"> <md:KeyDescriptor use=\"signing\"><ds:KeyInfo xmlns:ds=\"http://www.w3.org/2000/09/xmldsig#\"><ds:X509Data><ds:X509Certificate>MIIDpDCCAoygAwIBAgIGAX7/5qPhMA0GCSqGSIb3DQEBCwUAMIGSMQswCQYDVQQGEwJVUzETMBEGA1UECAwKQ2FsaWZvcm5pYTEWMBQGA1UEBwwNU2FuIEZyYW5jaXNjbzENMAsGA1UECgwET2t0YTEUMBIGA1UECwwLU1NPUHJvdmlkZXIxEzARBgNVBAMMCmRldi00NTg0MjExHDAaBgkqhkiG9w0BCQEWDWluZm9Ab2t0YS5jb20wHhcNMjIwMjE2MDAxOTEyWhcNMzIwMjE2MDAyMDEyWjCBkjELMAkGA1UEBhMCVVMxEzARBgNVBAgMCkNhbGlmb3JuaWExFjAUBgNVBAcMDVNhbiBGcmFuY2lzY28xDTALBgNVBAoMBE9rdGExFDASBgNVBAsMC1NTT1Byb3ZpZGVyMRMwEQYDVQQDDApkZXYtNDU4NDIxMRwwGgYJKoZIhvcNAQkBFg1pbmZvQG9rdGEuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAxrBl7GKz52cRpxF9xCsirnRuMxnhFBaUrsHqAQrLqWmdlpNYZTVg+T9iQ+aq/iE68L+BRZcZniKIvW58wqqS0ltXVvIkXuDSvnvnkkI5yMIVErR20K8jSOKQm1FmK+fgAJ4koshFiu9oLiqu0Ejc0DuL3/XRsb4RuxjktKTb1khgBBtb+7idEk0sFR0RPefAweXImJkDHDm7SxjDwGJUubbqpdTxasPr0W+AHI1VUzsUsTiHAoyb0XDkYqHfDzhj/ZdIEl4zHQ3bEZvlD984ztAnmX2SuFLLKfXeAAGHei8MMixJvwxYkkPeYZ/5h8WgBZPP4heS2CPjwYExt29L8QIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQARjJFz++a9Z5IQGFzsZMrX2EDR5ML4xxUiQkbhld1S1PljOLcYFARDmUC2YYHOueU4ee8Jid9nPGEUebV/4Jok+b+oQh+dWMgiWjSLI7h5q4OYZ3VJtdlVwgMFt2iz+/4yBKMUZ50g3Qgg36vE34us+eKitg759JgCNsibxn0qtJgSPm0sgP2L6yTaLnoEUbXBRxCwynTSkp9ZijZqEzbhN0e2dWv7Rx/nfpohpDP6vEiFImKFHpDSv3M/5de1ytQzPFrZBYt9WlzlYwE1aD9FHCxdd+rWgYMVVoRaRmndpV/Rq3QUuDuFJtaoX11bC7ExkOpg9KstZzA63i3VcfYv</ds:X509Certificate></ds:X509Data></ds:KeyInfo></md:KeyDescriptor><md:SingleSignOnService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect\" Location=\"https://test.com/sso\"/></md:IDPSSODescriptor></md:EntityDescriptor>"
  }
}
`, context)
}

func testAccIAMWorkforcePoolWorkforcePoolProvider_destroy(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_workforce_pool" "my_pool" {
  workforce_pool_id = "my-pool-%{random_suffix}"
  parent            = "organizations/%{org_id}"
  location          = "global"
}
`, context)
}
