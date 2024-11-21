// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package cloudbuildv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccCloudbuildv2Connection_GheCompleteConnection(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"region":        envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudbuildv2ConnectionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudbuildv2Connection_GheCompleteConnection(context),
			},
			{
				ResourceName:            "google_cloudbuildv2_connection.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "name"},
			},
		},
	})
}

func TestAccCloudbuildv2Connection_GheConnection(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"region":        envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudbuildv2ConnectionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudbuildv2Connection_GheConnection(context),
			},
			{
				ResourceName:            "google_cloudbuildv2_connection.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations"},
			},
			{
				Config: testAccCloudbuildv2Connection_GheConnectionUpdate0(context),
			},
			{
				ResourceName:            "google_cloudbuildv2_connection.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations"},
			},
		},
	})
}

func TestAccCloudbuildv2Connection_GhePrivConnection(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"region":        envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudbuildv2ConnectionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudbuildv2Connection_GhePrivConnection(context),
			},
			{
				ResourceName:            "google_cloudbuildv2_connection.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations"},
			},
		},
	})
}

func TestAccCloudbuildv2Connection_GhePrivUpdateConnection(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"region":        envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudbuildv2ConnectionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudbuildv2Connection_GhePrivUpdateConnection(context),
			},
			{
				ResourceName:            "google_cloudbuildv2_connection.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations"},
			},
			{
				Config: testAccCloudbuildv2Connection_GhePrivUpdateConnectionUpdate0(context),
			},
			{
				ResourceName:            "google_cloudbuildv2_connection.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations"},
			},
		},
	})
}

func TestAccCloudbuildv2Connection_GithubConnection(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"region":        envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudbuildv2ConnectionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudbuildv2Connection_GithubConnection(context),
			},
			{
				ResourceName:            "google_cloudbuildv2_connection.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations"},
			},
			{
				Config: testAccCloudbuildv2Connection_GithubConnectionUpdate0(context),
			},
			{
				ResourceName:            "google_cloudbuildv2_connection.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations"},
			},
		},
	})
}

func TestAccCloudbuildv2Connection_GitlabConnection(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudbuildv2ConnectionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudbuildv2Connection_GitlabConnection(context),
			},
			{
				ResourceName:            "google_cloudbuildv2_connection.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations"},
			},
		},
	})
}

func TestAccCloudbuildv2Connection_GleConnection(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudbuildv2ConnectionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudbuildv2Connection_GleConnection(context),
			},
			{
				ResourceName:            "google_cloudbuildv2_connection.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations"},
			},
			{
				Config: testAccCloudbuildv2Connection_GleConnectionUpdate0(context),
			},
			{
				ResourceName:            "google_cloudbuildv2_connection.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations"},
			},
		},
	})
}

func TestAccCloudbuildv2Connection_GleOldConnection(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudbuildv2ConnectionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudbuildv2Connection_GleOldConnection(context),
			},
			{
				ResourceName:            "google_cloudbuildv2_connection.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations"},
			},
			{
				Config: testAccCloudbuildv2Connection_GleOldConnectionUpdate0(context),
			},
			{
				ResourceName:            "google_cloudbuildv2_connection.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations"},
			},
		},
	})
}

func TestAccCloudbuildv2Connection_GlePrivConnection(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudbuildv2ConnectionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudbuildv2Connection_GlePrivConnection(context),
			},
			{
				ResourceName:            "google_cloudbuildv2_connection.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations"},
			},
		},
	})
}

func TestAccCloudbuildv2Connection_GlePrivUpdateConnection(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudbuildv2ConnectionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudbuildv2Connection_GleConnection(context),
			},
			{
				ResourceName:            "google_cloudbuildv2_connection.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations"},
			},
			{
				Config: testAccCloudbuildv2Connection_GlePrivConnection(context),
			},
			{
				ResourceName:            "google_cloudbuildv2_connection.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations"},
			},
		},
	})
}

func TestAccCloudbuildv2Connection_BbdcPrivConnection(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudbuildv2ConnectionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudbuildv2Connection_BbdcPrivConnection(context),
			},
			{
				ResourceName:            "google_cloudbuildv2_connection.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations"},
			},
		},
	})
}

func TestAccCloudbuildv2Connection_BbdcPrivUpdateConnection(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudbuildv2ConnectionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudbuildv2Connection_BbdcConnection(context),
			},
			{
				ResourceName:            "google_cloudbuildv2_connection.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations"},
			},
			{
				Config: testAccCloudbuildv2Connection_BbdcPrivConnection(context),
			},
			{
				ResourceName:            "google_cloudbuildv2_connection.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations"},
			},
		},
	})
}

func TestAccCloudbuildv2Connection_BbcConnection(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudbuildv2ConnectionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudbuildv2Connection_BbcConnection(context),
			},
			{
				ResourceName:            "google_cloudbuildv2_connection.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations"},
			},
		},
	})
}

func testAccCloudbuildv2Connection_GheCompleteConnection(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloudbuildv2_connection" "primary" {
  location = "%{region}"
  name     = "projects/%{project_name}/locations/%{region}/connections/tf-test-connection%{random_suffix}"

  github_enterprise_config {
    host_uri                      = "https://ghe.proctor-staging-test.com"
    app_id                        = 516
    app_installation_id           = 243
    app_slug                      = "myapp"
    private_key_secret_version    = "projects/gcb-terraform-creds/secrets/ghe-private-key/versions/latest"
    webhook_secret_secret_version = "projects/gcb-terraform-creds/secrets/ghe-webhook-secret/versions/latest"
  }

  project     = "%{project_name}"
  annotations = {}
}


`, context)
}

func testAccCloudbuildv2Connection_GheConnection(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloudbuildv2_connection" "primary" {
  location = "%{region}"
  name     = "tf-test-connection%{random_suffix}"

  github_enterprise_config {
    host_uri = "https://ghe.proctor-staging-test.com"
  }

  project     = "%{project_name}"
  annotations = {}
}


`, context)
}

func testAccCloudbuildv2Connection_GheConnectionUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloudbuildv2_connection" "primary" {
  location = "%{region}"
  name     = "tf-test-connection%{random_suffix}"

  github_enterprise_config {
    host_uri                      = "https://ghe.proctor-staging-test.com"
    app_id                        = 516
    app_installation_id           = 243
    app_slug                      = "myapp"
    private_key_secret_version    = "projects/gcb-terraform-creds/secrets/ghe-private-key/versions/latest"
    webhook_secret_secret_version = "projects/gcb-terraform-creds/secrets/ghe-webhook-secret/versions/latest"
  }

  project     = "%{project_name}"
  annotations = {}
}


`, context)
}

func testAccCloudbuildv2Connection_GhePrivConnection(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloudbuildv2_connection" "primary" {
  location = "%{region}"
  name     = "tf-test-connection%{random_suffix}"

  github_enterprise_config {
    host_uri = "https://ghe.proctor-private-ca.com"

    service_directory_config {
      service = "projects/gcb-terraform-creds/locations/%{region}/namespaces/myns/services/serv"
    }

    ssl_ca = "-----BEGIN CERTIFICATE-----\nMIIEXTCCA0WgAwIBAgIUANaBCc9j/xdKJHU0sgmv6yE2WCIwDQYJKoZIhvcNAQEL\nBQAwLDEUMBIGA1UEChMLUHJvY3RvciBFbmcxFDASBgNVBAMTC1Byb2N0b3ItZW5n\nMB4XDTIxMDcxNTIwMDcwMloXDTIyMDcxNTIwMDcwMVowADCCASIwDQYJKoZIhvcN\nAQEBBQADggEPADCCAQoCggEBAMVel7I88DkhwW445BNPBZvJNTV1AreHdz4um4U1\nop2+4L7JeNrUs5SRc0fzeOyOmA9ZzTDu9hBC7zj/sVNUy6cIQGCj32sr5SCAEIat\nnFZlzmVqJPT4J5NAaE37KO5347myTJEBrvpq8az4CtvX0yUzPK0gbUmaSaztVi4o\ndbJLKyv575xCLC/Hu6fIHBDH19eG1Ath9VpuAOkttRRoxu2VqijJZrGqaS+0o+OX\nrLi5HMtZbZjgQB4mc1g3ZDKX/gynxr+CDNaqNOqxuog33Tl5OcOk9DrR3MInaE7F\nyQFuH9mzF64AqOoTf7Tr/eAIz5XVt8K51nk+fSybEfKVwtMCAwEAAaOCAaEwggGd\nMA4GA1UdDwEB/wQEAwIFoDAMBgNVHRMBAf8EAjAAMB0GA1UdDgQWBBQU/9dYyqMz\nv9rOMwPZcoIRMDAQCjAfBgNVHSMEGDAWgBTkQGTiCkLCmv/Awxdz5TAVRmyFfDCB\njQYIKwYBBQUHAQEEgYAwfjB8BggrBgEFBQcwAoZwaHR0cDovL3ByaXZhdGVjYS1j\nb250ZW50LTYxYWEyYzA5LTAwMDAtMjJjMi05ZjYyLWQ0ZjU0N2Y4MDIwMC5zdG9y\nYWdlLmdvb2dsZWFwaXMuY29tLzQxNGU4ZTJjZjU2ZWEyYzQxNmM0L2NhLmNydDAo\nBgNVHREBAf8EHjAcghpnaGUucHJvY3Rvci1wcml2YXRlLWNhLmNvbTCBggYDVR0f\nBHsweTB3oHWgc4ZxaHR0cDovL3ByaXZhdGVjYS1jb250ZW50LTYxYWEyYzA5LTAw\nMDAtMjJjMi05ZjYyLWQ0ZjU0N2Y4MDIwMC5zdG9yYWdlLmdvb2dsZWFwaXMuY29t\nLzQxNGU4ZTJjZjU2ZWEyYzQxNmM0L2NybC5jcmwwDQYJKoZIhvcNAQELBQADggEB\nABo6BQLEZZ+YNiDuv2sRvcxSopQQb7fZjqIA9XOA35pNSKay2SncODnNvfsdRnOp\ncoy25sQSIzWyJ9zWl8DZ6evoOu5csZ2PoFqx5LsIq37w+ZcwD6DM8Zm7JqASxmxx\nGqTF0nHC4Aw8q8aJBeRD3PsSkfN5Q3DP3nTDnLyd0l+yPIkHUbZMoiFHX3BkhCng\nG96mYy/y3t16ghfV9lZkXpD/JK5aiN0bTHCDRc69owgfYiAcAqzBJ9gfZ90MBgzv\ngTTQel5dHg49SYXfnUpTy0HdQLEcoggOF8Q8V+xKdKa6eVbrvjJrkEJmvIQI5iCR\nhNvKR25mx8JUopqEXmONmqU=\n-----END CERTIFICATE-----\n\n-----BEGIN CERTIFICATE-----\nMIIDSDCCAjCgAwIBAgITMwWN+62nLcgyLa7p+jD1K90g6TANBgkqhkiG9w0BAQsF\nADAsMRQwEgYDVQQKEwtQcm9jdG9yIEVuZzEUMBIGA1UEAxMLUHJvY3Rvci1lbmcw\nHhcNMjEwNzEyMTM1OTQ0WhcNMzEwNzEwMTM1OTQzWjAsMRQwEgYDVQQKEwtQcm9j\ndG9yIEVuZzEUMBIGA1UEAxMLUHJvY3Rvci1lbmcwggEiMA0GCSqGSIb3DQEBAQUA\nA4IBDwAwggEKAoIBAQCYqJP5Qt90jIbld2dtuUV/zIkBFsTe4fapJfhBji03xBpN\nO1Yxj/jPSZ67Kdeoy0lEwvc2hL5FQGhIjLMR0mzOyN4fk/DZiA/4tAVi7hJyqpUC\n71JSwp7MwXL1b26CSE1MhcoCqA/E4iZxfJfF/ef4lhmC24UEmu8FEbldoy+6OysB\nRu7dGDwicW5F9h7eSkpGAsCRdJHh65iUx/IH0C4Ux2UZRDZdj6wVbuVu9tb938xF\nyRuVClONoLSn/lwdzeV7hQmBSm8qmfgbNPbYRaNLz3hOpsT+27aDQp2/pxue8hFJ\nd7We3+Lr5O4IL45PBwhVEAiFZqde6d4qViNEB2qTAgMBAAGjYzBhMA4GA1UdDwEB\n/wQEAwIBBjAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBTkQGTiCkLCmv/Awxdz\n5TAVRmyFfDAfBgNVHSMEGDAWgBTkQGTiCkLCmv/Awxdz5TAVRmyFfDANBgkqhkiG\n9w0BAQsFAAOCAQEAfy5BJsWdx0oWWi7SFg9MbryWjBVPJl93UqACgG0Cgh813O/x\nlDZQhGO/ZFVhHz/WgooE/HgVNoVJTubKLLzz+zCkOB0wa3GMqJDyFjhFmUtd/3VM\nZh0ZQ+JWYsAiZW4VITj5xEn/d/B3xCFWGC1vhvhptEJ8Fo2cE1yM2pzk08NqFWoY\n4FaH0sbxWgyCKwTmtcYDbnx4FYuddryGCIxbYizqUK1dr4DGKeHonhm/d234Ew3x\n3vIBPoHMOfBec/coP1xAf5o+F+MRMO/sQ3tTGgyOH18lwsHo9SmXCrmOwVQPKrEw\nm+A+5TjXLmenyaBhqXa0vkAZYJhWdROhWC0VTA==\n-----END CERTIFICATE-----\n"
  }

  project     = "%{project_name}"
  annotations = {}
}


`, context)
}

func testAccCloudbuildv2Connection_GhePrivUpdateConnection(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloudbuildv2_connection" "primary" {
  location = "%{region}"
  name     = "tf-test-connection%{random_suffix}"

  github_enterprise_config {
    host_uri = "https://ghe.proctor-staging-test.com"
  }

  project     = "%{project_name}"
  annotations = {}
}


`, context)
}

func testAccCloudbuildv2Connection_GhePrivUpdateConnectionUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloudbuildv2_connection" "primary" {
  location = "%{region}"
  name     = "tf-test-connection%{random_suffix}"

  github_enterprise_config {
    host_uri = "https://ghe.proctor-private-ca.com"

    service_directory_config {
      service = "projects/gcb-terraform-creds/locations/%{region}/namespaces/myns/services/serv"
    }

    ssl_ca = "-----BEGIN CERTIFICATE-----\nMIIEXTCCA0WgAwIBAgIUANaBCc9j/xdKJHU0sgmv6yE2WCIwDQYJKoZIhvcNAQEL\nBQAwLDEUMBIGA1UEChMLUHJvY3RvciBFbmcxFDASBgNVBAMTC1Byb2N0b3ItZW5n\nMB4XDTIxMDcxNTIwMDcwMloXDTIyMDcxNTIwMDcwMVowADCCASIwDQYJKoZIhvcN\nAQEBBQADggEPADCCAQoCggEBAMVel7I88DkhwW445BNPBZvJNTV1AreHdz4um4U1\nop2+4L7JeNrUs5SRc0fzeOyOmA9ZzTDu9hBC7zj/sVNUy6cIQGCj32sr5SCAEIat\nnFZlzmVqJPT4J5NAaE37KO5347myTJEBrvpq8az4CtvX0yUzPK0gbUmaSaztVi4o\ndbJLKyv575xCLC/Hu6fIHBDH19eG1Ath9VpuAOkttRRoxu2VqijJZrGqaS+0o+OX\nrLi5HMtZbZjgQB4mc1g3ZDKX/gynxr+CDNaqNOqxuog33Tl5OcOk9DrR3MInaE7F\nyQFuH9mzF64AqOoTf7Tr/eAIz5XVt8K51nk+fSybEfKVwtMCAwEAAaOCAaEwggGd\nMA4GA1UdDwEB/wQEAwIFoDAMBgNVHRMBAf8EAjAAMB0GA1UdDgQWBBQU/9dYyqMz\nv9rOMwPZcoIRMDAQCjAfBgNVHSMEGDAWgBTkQGTiCkLCmv/Awxdz5TAVRmyFfDCB\njQYIKwYBBQUHAQEEgYAwfjB8BggrBgEFBQcwAoZwaHR0cDovL3ByaXZhdGVjYS1j\nb250ZW50LTYxYWEyYzA5LTAwMDAtMjJjMi05ZjYyLWQ0ZjU0N2Y4MDIwMC5zdG9y\nYWdlLmdvb2dsZWFwaXMuY29tLzQxNGU4ZTJjZjU2ZWEyYzQxNmM0L2NhLmNydDAo\nBgNVHREBAf8EHjAcghpnaGUucHJvY3Rvci1wcml2YXRlLWNhLmNvbTCBggYDVR0f\nBHsweTB3oHWgc4ZxaHR0cDovL3ByaXZhdGVjYS1jb250ZW50LTYxYWEyYzA5LTAw\nMDAtMjJjMi05ZjYyLWQ0ZjU0N2Y4MDIwMC5zdG9yYWdlLmdvb2dsZWFwaXMuY29t\nLzQxNGU4ZTJjZjU2ZWEyYzQxNmM0L2NybC5jcmwwDQYJKoZIhvcNAQELBQADggEB\nABo6BQLEZZ+YNiDuv2sRvcxSopQQb7fZjqIA9XOA35pNSKay2SncODnNvfsdRnOp\ncoy25sQSIzWyJ9zWl8DZ6evoOu5csZ2PoFqx5LsIq37w+ZcwD6DM8Zm7JqASxmxx\nGqTF0nHC4Aw8q8aJBeRD3PsSkfN5Q3DP3nTDnLyd0l+yPIkHUbZMoiFHX3BkhCng\nG96mYy/y3t16ghfV9lZkXpD/JK5aiN0bTHCDRc69owgfYiAcAqzBJ9gfZ90MBgzv\ngTTQel5dHg49SYXfnUpTy0HdQLEcoggOF8Q8V+xKdKa6eVbrvjJrkEJmvIQI5iCR\nhNvKR25mx8JUopqEXmONmqU=\n-----END CERTIFICATE-----\n\n-----BEGIN CERTIFICATE-----\nMIIDSDCCAjCgAwIBAgITMwWN+62nLcgyLa7p+jD1K90g6TANBgkqhkiG9w0BAQsF\nADAsMRQwEgYDVQQKEwtQcm9jdG9yIEVuZzEUMBIGA1UEAxMLUHJvY3Rvci1lbmcw\nHhcNMjEwNzEyMTM1OTQ0WhcNMzEwNzEwMTM1OTQzWjAsMRQwEgYDVQQKEwtQcm9j\ndG9yIEVuZzEUMBIGA1UEAxMLUHJvY3Rvci1lbmcwggEiMA0GCSqGSIb3DQEBAQUA\nA4IBDwAwggEKAoIBAQCYqJP5Qt90jIbld2dtuUV/zIkBFsTe4fapJfhBji03xBpN\nO1Yxj/jPSZ67Kdeoy0lEwvc2hL5FQGhIjLMR0mzOyN4fk/DZiA/4tAVi7hJyqpUC\n71JSwp7MwXL1b26CSE1MhcoCqA/E4iZxfJfF/ef4lhmC24UEmu8FEbldoy+6OysB\nRu7dGDwicW5F9h7eSkpGAsCRdJHh65iUx/IH0C4Ux2UZRDZdj6wVbuVu9tb938xF\nyRuVClONoLSn/lwdzeV7hQmBSm8qmfgbNPbYRaNLz3hOpsT+27aDQp2/pxue8hFJ\nd7We3+Lr5O4IL45PBwhVEAiFZqde6d4qViNEB2qTAgMBAAGjYzBhMA4GA1UdDwEB\n/wQEAwIBBjAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBTkQGTiCkLCmv/Awxdz\n5TAVRmyFfDAfBgNVHSMEGDAWgBTkQGTiCkLCmv/Awxdz5TAVRmyFfDANBgkqhkiG\n9w0BAQsFAAOCAQEAfy5BJsWdx0oWWi7SFg9MbryWjBVPJl93UqACgG0Cgh813O/x\nlDZQhGO/ZFVhHz/WgooE/HgVNoVJTubKLLzz+zCkOB0wa3GMqJDyFjhFmUtd/3VM\nZh0ZQ+JWYsAiZW4VITj5xEn/d/B3xCFWGC1vhvhptEJ8Fo2cE1yM2pzk08NqFWoY\n4FaH0sbxWgyCKwTmtcYDbnx4FYuddryGCIxbYizqUK1dr4DGKeHonhm/d234Ew3x\n3vIBPoHMOfBec/coP1xAf5o+F+MRMO/sQ3tTGgyOH18lwsHo9SmXCrmOwVQPKrEw\nm+A+5TjXLmenyaBhqXa0vkAZYJhWdROhWC0VTA==\n-----END CERTIFICATE-----\n"
  }

  project     = "%{project_name}"
  annotations = {}
}


`, context)
}

func testAccCloudbuildv2Connection_GithubConnection(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloudbuildv2_connection" "primary" {
  location = "%{region}"
  name     = "tf-test-connection%{random_suffix}"
  disabled = true

  github_config {
    app_installation_id = 0

    authorizer_credential {
      oauth_token_secret_version = "projects/gcb-terraform-creds/secrets/github-pat/versions/1"
    }
  }

  project = "%{project_name}"

  annotations = {
    somekey = "somevalue"
  }
}


`, context)
}

func testAccCloudbuildv2Connection_GithubConnectionUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloudbuildv2_connection" "primary" {
  location = "%{region}"
  name     = "tf-test-connection%{random_suffix}"
  disabled = false

  github_config {
    app_installation_id = 31300675

    authorizer_credential {
      oauth_token_secret_version = "projects/gcb-terraform-creds/secrets/github-pat/versions/latest"
    }
  }

  project = "%{project_name}"

  annotations = {
    otherkey = "othervalue"

    somekey = "somevalue"
  }
}


`, context)
}

func testAccCloudbuildv2Connection_GitlabConnection(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloudbuildv2_connection" "primary" {
  location = "us-west1"
  name     = "tf-test-connection%{random_suffix}"

  gitlab_config {
    authorizer_credential {
      user_token_secret_version = "projects/407304063574/secrets/gitlab-api-pat/versions/latest"
    }

    read_authorizer_credential {
      user_token_secret_version = "projects/407304063574/secrets/gitlab-read-pat/versions/latest"
    }

    webhook_secret_secret_version = "projects/407304063574/secrets/gle-webhook-secret/versions/latest"
  }

  project     = "%{project_name}"
  annotations = {}
}


`, context)
}

func testAccCloudbuildv2Connection_GleConnection(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloudbuildv2_connection" "primary" {
  location = "us-west1"
  name     = "tf-test-connection%{random_suffix}"

  gitlab_config {
    authorizer_credential {
      user_token_secret_version = "projects/407304063574/secrets/gle-api-token/versions/latest"
    }

    read_authorizer_credential {
      user_token_secret_version = "projects/407304063574/secrets/gle-read-token/versions/latest"
    }

    webhook_secret_secret_version = "projects/407304063574/secrets/gle-webhook-secret/versions/latest"
    host_uri                      = "https://gle-us-central1.gcb-test.com"
  }

  project     = "%{project_name}"
  annotations = {}
}


`, context)
}

func testAccCloudbuildv2Connection_GleConnectionUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloudbuildv2_connection" "primary" {
  location = "us-west1"
  name     = "tf-test-connection%{random_suffix}"

  gitlab_config {
    authorizer_credential {
      user_token_secret_version = "projects/407304063574/secrets/gle-old-api-token/versions/2"
    }

    read_authorizer_credential {
      user_token_secret_version = "projects/407304063574/secrets/gle-old-read-token/versions/3"
    }

    webhook_secret_secret_version = "projects/407304063574/secrets/gle-webhook-secret/versions/latest"
    host_uri                      = "https://gle-old.gcb-test.com"
  }

  project     = "%{project_name}"
  annotations = {}
}


`, context)
}

func testAccCloudbuildv2Connection_GleOldConnection(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloudbuildv2_connection" "primary" {
  location = "us-west1"
  name     = "tf-test-connection%{random_suffix}"

  gitlab_config {
    authorizer_credential {
      user_token_secret_version = "projects/407304063574/secrets/gle-old-api-token/versions/2"
    }

    read_authorizer_credential {
      user_token_secret_version = "projects/407304063574/secrets/gle-old-read-token/versions/3"
    }

    webhook_secret_secret_version = "projects/407304063574/secrets/gle-webhook-secret/versions/latest"
    host_uri                      = "https://gle-old.gcb-test.com"
  }

  project     = "%{project_name}"
  annotations = {}
}


`, context)
}

func testAccCloudbuildv2Connection_GleOldConnectionUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloudbuildv2_connection" "primary" {
  location = "us-west1"
  name     = "tf-test-connection%{random_suffix}"

  gitlab_config {
    authorizer_credential {
      user_token_secret_version = "projects/407304063574/secrets/gle-api-token/versions/latest"
    }

    read_authorizer_credential {
      user_token_secret_version = "projects/407304063574/secrets/gle-read-token/versions/latest"
    }

    webhook_secret_secret_version = "projects/407304063574/secrets/gle-webhook-secret/versions/latest"
    host_uri                      = "https://gle-us-central1.gcb-test.com"
  }

  project     = "%{project_name}"
  annotations = {}
}


`, context)
}

func testAccCloudbuildv2Connection_GlePrivConnection(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloudbuildv2_connection" "primary" {
  location = "us-west1"
  name     = "tf-test-connection%{random_suffix}"

  gitlab_config {
    authorizer_credential {
      user_token_secret_version = "projects/407304063574/secrets/gle-private-api/versions/latest"
    }

    read_authorizer_credential {
      user_token_secret_version = "projects/407304063574/secrets/gle-private-read-token/versions/latest"
    }

    webhook_secret_secret_version = "projects/407304063574/secrets/gle-webhook-secret/versions/latest"
    host_uri                      = "https://gle-us.gle-us-private.com"

    service_directory_config {
      service = "projects/407304063574/locations/us-west1/namespaces/private-conn/services/gitlab-private"
    }

	ssl_ca = "-----BEGIN CERTIFICATE-----\nMIIFbjCCA1agAwIBAgIUH+nsWsqagMW9Ld8E9J71yPLPpD8wDQYJKoZIhvcNAQEL\nBQAwJDEiMCAGA1UEAwwZZ2xlLXVzLmdsZS11cy1wcml2YXRlLmNvbTAeFw0yNDEw\nMzExNjQzMjBaFw0zNDEwMjkxNjQzMjBaMCQxIjAgBgNVBAMMGWdsZS11cy5nbGUt\ndXMtcHJpdmF0ZS5jb20wggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQDL\n+dUU8MHo+Eskx4SSnI1thRIiljgsyJSzSOplaD4lqahFnrG0cB0ovKpyRL4A+0wM\nzVW7W1Pfi8DiEOxxfNo7pEj+0zrzJHHqnzW9kApIlRmO1TBBJ7i9HaVamJ1Od01b\n2WI/pnKFEvNfLQSDQBulhkHZ2McyauDhb1DxefKnVX8ac6qhxtc4IzrezIQuJ18N\nDPtNLUDD4rtU4mIX4lx1yBIplrgypAo1HDbJOwW8OR76MtjAY7ek3K6UCyduQtwy\nmfZ23b3Eg69W10HVMVTy6m5NaGKi/TWy2MJ71hKUQ1+tWIPH5SL7FzYPKL4XXw5W\np61LhIiBAd2tgD41b2cQxhUbVifc1qHtnwNz/tE77M9ySH37rEUIlExzr3D3JV+f\nXjXEXUr9as8GRnS5zhD/opKe7wKbwpYMHhylK1h6XH/sBO7dBajf5xCvpZZBDzrK\nkpTqwHspT7p40WF9d8odjEk/xZKn5LdcDG2I+4U7SVS5e8ud41HUQxJwQx56lKfh\n2WB+zs7nSyMfspTj4doY1OADEC1VQCyGrwlbclKTKmUWrgwQdm38KxDzW5Juyjmm\nzvfsWIlSMdnes0qVVo38N3Jz8/MHCLD79R2veWgA2fbqS3+4h2dRkR7htjaVlJMJ\nt7SwFiG39ic3OZpo+wTkaHlG4CBnbFDueUsOW2wEpQIDAQABo4GXMIGUMB0GA1Ud\nDgQWBBTExgzH2gz9+rJHvlTFPO0AvG88azAfBgNVHSMEGDAWgBTExgzH2gz9+rJH\nvlTFPO0AvG88azAPBgNVHRMBAf8EBTADAQH/MEEGA1UdEQQ6MDiCGWdsZS11cy5n\nbGUtdXMtcHJpdmF0ZS5jb22CGyouZ2xlLXVzLmdsZS11cy1wcml2YXRlLmNvbTAN\nBgkqhkiG9w0BAQsFAAOCAgEAjkd1ZNoekoWrmozD+Ta1OM0zWhv04eqhP8aYzhbd\nXRS+GyF6ifMwfWg9HogkH22ZPT5GszaL5DacSyOUqZgJ905Q6g1EFPnaKmFVHHeC\nzZAhg5oedAzcakZpYwZDSiLuPgsQfwgRnqWIYR8JcIM5bKRZNGyOg8eZ8cKu23A2\nPavL4B3Ra1l93KllKm21rigIhLPIPLoEyxEg9c9oTJF92r0+aRdf2Ln853260Fqf\ncEUWoXhqMGvDv/YEbqDjGQ/Kh7ZWdlIWhcKFOA0gluF7oExjt/MgSitukgg3aaic\n/eXXOrZDNYH7Ve610NUuNlhub1M47Tp7EgjUJVWlsKK84T8ZcZq7Hn4BzioUr95d\nHao6u19HWA/ISM8bwzHaYxscFI4u6phEL0HJzLf4EysEmS0rAnLxyol0apNx6znR\nhXsqxnSexKhXoLqnK1Vuhcg8DsvobXHqg68EGZ7BZ3ycPYaHSWU8Xh3l1gtYkcQ6\nzxXsKIijlpVKuYJvGA3EOMoZu6+2MYF8Tgp3N4sKMvPhqBhsmgxOYF5OkAbGlsUP\nyCYWFDBFHmbhvUu5JpbKuID2CPkBi16EetemvMQ9PGlLq/0fO/BBNkn6TYn9Kvg8\nAyvuONz54uFEAIKPCcZIosa3ml+5/pt+tBhtVzHA6vMxn18IYaNpuTwSxi/+M10K\nRjw=\n-----END CERTIFICATE-----\n"
  }

  project     = "%{project_name}"
  annotations = {}
}


`, context)
}

func testAccCloudbuildv2Connection_BbdcConnection(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloudbuildv2_connection" "primary" {
  location = "us-west1"
  name     = "tf-test-connection%{random_suffix}"

  bitbucket_data_center_config {
    authorizer_credential {
      user_token_secret_version = "projects/407304063574/secrets/bbdc-api-token/versions/latest"
    }

    read_authorizer_credential {
      user_token_secret_version = "projects/407304063574/secrets/bbdc-read-token/versions/latest"
    }

    webhook_secret_secret_version = "projects/407304063574/secrets/bbdc-webhook-secret/versions/latest"
    host_uri                      = "https://bitbucket-us-central.gcb-test.com"
  }

  project     = "%{project_name}"
  annotations = {}
}


`, context)
}

func testAccCloudbuildv2Connection_BbdcPrivConnection(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloudbuildv2_connection" "primary" {
  location = "us-west1"
  name     = "tf-test-connection%{random_suffix}"

  bitbucket_data_center_config {
    authorizer_credential {
      user_token_secret_version = "projects/407304063574/secrets/private-bbdc-api-token/versions/1"
    }

    read_authorizer_credential {
      user_token_secret_version = "projects/407304063574/secrets/private-bbdc-read-token/versions/1"
    }

    webhook_secret_secret_version = "projects/407304063574/secrets/bbdc-webhook-secret/versions/latest"
    host_uri                      = "https://private-bitbucket.proctor-test.com"

  service_directory_config {
    service = "projects/407304063574/locations/us-west1/namespaces/private-conn/services/private-bitbucket"
  }

	ssl_ca = "-----BEGIN CERTIFICATE-----\nMIIDjDCCAnSgAwIBAgIUBh5+3oeT1vmUSS5rSNaFfy6igSAwDQYJKoZIhvcNAQEL\nBQAwVzELMAkGA1UEBhMCVVMxGzAZBgNVBAoMEkdvb2dsZSBDbG91ZCBCdWlsZDEr\nMCkGA1UEAwwicHJpdmF0ZS1iaXRidWNrZXQucHJvY3Rvci10ZXN0LmNvbTAeFw0y\nMzEyMTIyMzI5NTlaFw0yNDEyMTEyMzI5NTlaMFcxCzAJBgNVBAYTAlVTMRswGQYD\nVQQKDBJHb29nbGUgQ2xvdWQgQnVpbGQxKzApBgNVBAMMInByaXZhdGUtYml0YnVj\na2V0LnByb2N0b3ItdGVzdC5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK\nAoIBAQCNfMx4ImGD4imZR64RbmRtpUNmypDokx2/S9kgobmyNvBWeSgRVhOGHGbU\nUgyvENcEg803K8unwF2jF6sdGrocRnIdPpr2tUoViOM2Ss6ds+TD8a2kqBA6+hmQ\nOMJiEIirpGT3Mw1pYTpuLisfIeeuuYssoS5k18kFLZ+Mk6MUSAHCgC8EowUZLGBZ\nagh9OhrjpMSXyidv+2d7FKTh/k3BWffVkDXehjvWjcr47hSvQwqW5m773ewCq0uD\nwxUgO6MAAAxLJz15cjhfvk4ishgSqcp49IZrx+xsNCLbHjPVyGkrL2OhgFaGsQS/\nq6GkXYfJ1sJYrf5Xm1EXbZlQZzJPAgMBAAGjUDBOMC0GA1UdEQQmMCSCInByaXZh\ndGUtYml0YnVja2V0LnByb2N0b3ItdGVzdC5jb20wHQYDVR0OBBYEFISmuuTpHKMB\n+m1h62gEqg1ovC86MA0GCSqGSIb3DQEBCwUAA4IBAQAwIwR6pIum9EZyLtC438Q1\nEgH3SKqbdyMFCkFSBvr4WfFU6ja1pn5ZxzJWt5TRFlI9GMy7BupQrxJGebOiFuUC\noNJpc4QDt9a0/GKh48DGF7uKo9XK33p0v1ahq3ewNT/CUnHewQNX7aXXP1/rL+br\nZPA20XWURUTviMik7DdhaXKQv76K9coI3H74heeBUp+OHKgUkqA3D1QIGNRGOKos\n4z6MyBWVpMUIeJQGtIQBd9CY1hBN231iG1+hdOlOMwgyNVK2GS738r+HbngFo9v4\nh2I1HMUHVcHiPQLqwZ2/OTmTmF1aWCUbhnAvoisu20rHVcGnVIOqMrHYFzdGr3ZQ\n-----END CERTIFICATE-----\n"  
  }

  project     = "%{project_name}"
  annotations = {}
}


`, context)
}

func testAccCloudbuildv2Connection_BbcConnection(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloudbuildv2_connection" "primary" {
  location = "us-west1"
  name     = "tf-test-connection%{random_suffix}"

  bitbucket_cloud_config {
  workspace = "proctor-test"
    authorizer_credential {
      user_token_secret_version = "projects/407304063574/secrets/bbc-api-token/versions/latest"
    }

    read_authorizer_credential {
      user_token_secret_version = "projects/407304063574/secrets/bbc-read-token/versions/latest"
    }

    webhook_secret_secret_version = "projects/407304063574/secrets/bbdc-webhook-secret/versions/latest"
  }

  project     = "%{project_name}"
  annotations = {}
}


`, context)
}
