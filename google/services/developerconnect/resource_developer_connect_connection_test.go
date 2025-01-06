// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package developerconnect_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDeveloperConnectConnection_developerConnectConnectionGithubUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDeveloperConnectConnection_Github(context),
			},
			{
				ResourceName:            "google_developer_connect_connection.my-connection",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "connection_id", "labels", "location", "terraform_labels"},
			},
			{
				Config: testAccDeveloperConnectConnection_GithubUpdate(context),
			},
			{
				ResourceName:            "google_developer_connect_connection.my-connection",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"connection_id", "location", "terraform_labels"},
			},
		},
	})
}

func testAccDeveloperConnectConnection_Github(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_developer_connect_connection" "my-connection" {
  location = "us-central1"
  connection_id = "tf-test-tf-test-connection%{random_suffix}"

  github_config {
    github_app = "DEVELOPER_CONNECT"

    authorizer_credential {
      oauth_token_secret_version = "projects/devconnect-terraform-creds/secrets/tf-test-do-not-change-github-oauthtoken-e0b9e7/versions/1"
    }
  }
}
`, context)
}

func testAccDeveloperConnectConnection_GithubUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_developer_connect_connection" "my-connection" {
  location = "us-central1"
  connection_id = "tf-test-tf-test-connection%{random_suffix}"
  annotations = {}
  labels = {}
  
  crypto_key_config {
    key_reference = "projects/devconnect-terraform-creds/locations/us-central1/keyRings/tf-keyring/cryptoKeys/tf-crypto-key"
  }

  github_config {
    github_app = "DEVELOPER_CONNECT"
    app_installation_id = 49439208

    authorizer_credential {
      oauth_token_secret_version = "projects/devconnect-terraform-creds/secrets/tf-test-do-not-change-github-oauthtoken-e0b9e7/versions/1"
    }
  }
}
`, context)
}

func TestAccDeveloperConnectConnection_developerConnectConnectionGithubEnterpriseUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDeveloperConnectConnection_GithubEnterprise(context),
			},
			{
				ResourceName:            "google_developer_connect_connection.my-connection",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"connection_id", "labels", "terraform_labels"},
			},
			{
				Config: testAccDeveloperConnectConnection_GithubEnterpriseUpdate(context),
			},
			{
				ResourceName:            "google_developer_connect_connection.my-connection",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"connection_id", "location", "terraform_labels"},
			},
		},
	})
}

func testAccDeveloperConnectConnection_GithubEnterprise(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_developer_connect_connection" "my-connection" {
  location = "us-central1"
  connection_id = "tf-test-tf-test-connection%{random_suffix}"

  github_enterprise_config {
    host_uri = "https://ghe.proctor-staging-test.com"
    app_id = 864434
    private_key_secret_version = "projects/devconnect-terraform-creds/secrets/tf-test-ghe-do-not-change-ghe-private-key-f522d2/versions/latest"
    webhook_secret_secret_version = "projects/devconnect-terraform-creds/secrets/tf-test-ghe-do-not-change-ghe-webhook-secret-3c806f/versions/latest"
  }
}
`, context)
}

func testAccDeveloperConnectConnection_GithubEnterpriseUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_developer_connect_connection" "my-connection" {
  location = "us-central1"
  connection_id = "tf-test-tf-test-connection%{random_suffix}"
  annotations = {}
  labels = {}
  
  crypto_key_config {
    key_reference = "projects/devconnect-terraform-creds/locations/us-central1/keyRings/tf-keyring/cryptoKeys/tf-crypto-key"
  }

  github_enterprise_config {
    host_uri = "https://ghe-asia.proctor-staging-test.com"
    app_id = 866372
    private_key_secret_version = "projects/devconnect-terraform-creds/secrets/ghe-private-key-update/versions/latest"
    webhook_secret_secret_version = "projects/devconnect-terraform-creds/secrets/ghe-webhook-secret-update/versions/latest"
    app_installation_id = 808867
  }
}
`, context)
}

func TestAccDeveloperConnectConnection_GhePrivConnection(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDeveloperConnectConnection_GhePrivConnection(context),
			},
			{
				ResourceName:            "google_developer_connect_connection.my-connection",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"connection_id", "location", "terraform_labels"},
			},
		},
	})
}

func testAccDeveloperConnectConnection_GhePrivConnection(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_developer_connect_connection" "my-connection" {
  location = "us-central1"
  connection_id = "tf-test-tf-test-connection%{random_suffix}"
  annotations = {}
  labels = {}

  github_enterprise_config {
    host_uri = "https://ghe.proctor-private-ca.com"
    app_id = 26
    private_key_secret_version = "projects/devconnect-terraform-creds/secrets/ghe-priv-private-key/versions/latest"
    webhook_secret_secret_version = "projects/devconnect-terraform-creds/secrets/tf-test-ghe-do-not-change-ghe-webhook-secret-3c806f/versions/latest"
    app_installation_id = 24
       
    ssl_ca_certificate = "-----BEGIN CERTIFICATE-----\nMIIEXTCCA0WgAwIBAgIUANaBCc9j/xdKJHU0sgmv6yE2WCIwDQYJKoZIhvcNAQEL\nBQAwLDEUMBIGA1UEChMLUHJvY3RvciBFbmcxFDASBgNVBAMTC1Byb2N0b3ItZW5n\nMB4XDTIxMDcxNTIwMDcwMloXDTIyMDcxNTIwMDcwMVowADCCASIwDQYJKoZIhvcN\nAQEBBQADggEPADCCAQoCggEBAMVel7I88DkhwW445BNPBZvJNTV1AreHdz4um4U1\nop2+4L7JeNrUs5SRc0fzeOyOmA9ZzTDu9hBC7zj/sVNUy6cIQGCj32sr5SCAEIat\nnFZlzmVqJPT4J5NAaE37KO5347myTJEBrvpq8az4CtvX0yUzPK0gbUmaSaztVi4o\ndbJLKyv575xCLC/Hu6fIHBDH19eG1Ath9VpuAOkttRRoxu2VqijJZrGqaS+0o+OX\nrLi5HMtZbZjgQB4mc1g3ZDKX/gynxr+CDNaqNOqxuog33Tl5OcOk9DrR3MInaE7F\nyQFuH9mzF64AqOoTf7Tr/eAIz5XVt8K51nk+fSybEfKVwtMCAwEAAaOCAaEwggGd\nMA4GA1UdDwEB/wQEAwIFoDAMBgNVHRMBAf8EAjAAMB0GA1UdDgQWBBQU/9dYyqMz\nv9rOMwPZcoIRMDAQCjAfBgNVHSMEGDAWgBTkQGTiCkLCmv/Awxdz5TAVRmyFfDCB\njQYIKwYBBQUHAQEEgYAwfjB8BggrBgEFBQcwAoZwaHR0cDovL3ByaXZhdGVjYS1j\nb250ZW50LTYxYWEyYzA5LTAwMDAtMjJjMi05ZjYyLWQ0ZjU0N2Y4MDIwMC5zdG9y\nYWdlLmdvb2dsZWFwaXMuY29tLzQxNGU4ZTJjZjU2ZWEyYzQxNmM0L2NhLmNydDAo\nBgNVHREBAf8EHjAcghpnaGUucHJvY3Rvci1wcml2YXRlLWNhLmNvbTCBggYDVR0f\nBHsweTB3oHWgc4ZxaHR0cDovL3ByaXZhdGVjYS1jb250ZW50LTYxYWEyYzA5LTAw\nMDAtMjJjMi05ZjYyLWQ0ZjU0N2Y4MDIwMC5zdG9yYWdlLmdvb2dsZWFwaXMuY29t\nLzQxNGU4ZTJjZjU2ZWEyYzQxNmM0L2NybC5jcmwwDQYJKoZIhvcNAQELBQADggEB\nABo6BQLEZZ+YNiDuv2sRvcxSopQQb7fZjqIA9XOA35pNSKay2SncODnNvfsdRnOp\ncoy25sQSIzWyJ9zWl8DZ6evoOu5csZ2PoFqx5LsIq37w+ZcwD6DM8Zm7JqASxmxx\nGqTF0nHC4Aw8q8aJBeRD3PsSkfN5Q3DP3nTDnLyd0l+yPIkHUbZMoiFHX3BkhCng\nG96mYy/y3t16ghfV9lZkXpD/JK5aiN0bTHCDRc69owgfYiAcAqzBJ9gfZ90MBgzv\ngTTQel5dHg49SYXfnUpTy0HdQLEcoggOF8Q8V+xKdKa6eVbrvjJrkEJmvIQI5iCR\nhNvKR25mx8JUopqEXmONmqU=\n-----END CERTIFICATE-----\n\n-----BEGIN CERTIFICATE-----\nMIIDSDCCAjCgAwIBAgITMwWN+62nLcgyLa7p+jD1K90g6TANBgkqhkiG9w0BAQsF\nADAsMRQwEgYDVQQKEwtQcm9jdG9yIEVuZzEUMBIGA1UEAxMLUHJvY3Rvci1lbmcw\nHhcNMjEwNzEyMTM1OTQ0WhcNMzEwNzEwMTM1OTQzWjAsMRQwEgYDVQQKEwtQcm9j\ndG9yIEVuZzEUMBIGA1UEAxMLUHJvY3Rvci1lbmcwggEiMA0GCSqGSIb3DQEBAQUA\nA4IBDwAwggEKAoIBAQCYqJP5Qt90jIbld2dtuUV/zIkBFsTe4fapJfhBji03xBpN\nO1Yxj/jPSZ67Kdeoy0lEwvc2hL5FQGhIjLMR0mzOyN4fk/DZiA/4tAVi7hJyqpUC\n71JSwp7MwXL1b26CSE1MhcoCqA/E4iZxfJfF/ef4lhmC24UEmu8FEbldoy+6OysB\nRu7dGDwicW5F9h7eSkpGAsCRdJHh65iUx/IH0C4Ux2UZRDZdj6wVbuVu9tb938xF\nyRuVClONoLSn/lwdzeV7hQmBSm8qmfgbNPbYRaNLz3hOpsT+27aDQp2/pxue8hFJ\nd7We3+Lr5O4IL45PBwhVEAiFZqde6d4qViNEB2qTAgMBAAGjYzBhMA4GA1UdDwEB\n/wQEAwIBBjAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBTkQGTiCkLCmv/Awxdz\n5TAVRmyFfDAfBgNVHSMEGDAWgBTkQGTiCkLCmv/Awxdz5TAVRmyFfDANBgkqhkiG\n9w0BAQsFAAOCAQEAfy5BJsWdx0oWWi7SFg9MbryWjBVPJl93UqACgG0Cgh813O/x\nlDZQhGO/ZFVhHz/WgooE/HgVNoVJTubKLLzz+zCkOB0wa3GMqJDyFjhFmUtd/3VM\nZh0ZQ+JWYsAiZW4VITj5xEn/d/B3xCFWGC1vhvhptEJ8Fo2cE1yM2pzk08NqFWoY\n4FaH0sbxWgyCKwTmtcYDbnx4FYuddryGCIxbYizqUK1dr4DGKeHonhm/d234Ew3x\n3vIBPoHMOfBec/coP1xAf5o+F+MRMO/sQ3tTGgyOH18lwsHo9SmXCrmOwVQPKrEw\nm+A+5TjXLmenyaBhqXa0vkAZYJhWdROhWC0VTA==\n-----END CERTIFICATE-----\n"
 
    service_directory_config {
      service = "projects/devconnect-terraform-creds/locations/us-central1/namespaces/my-namespace/services/terraform-github"
    }
  }

}
`, context)
}

func TestAccDeveloperConnectConnection_developerConnectConnectionGitlabUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDeveloperConnectConnection_Gitlab(context),
			},
			{
				ResourceName:            "google_developer_connect_connection.my-connection",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"connection_id", "location", "terraform_labels"},
			},
			{
				Config: testAccDeveloperConnectConnection_GitlabUpdate(context),
			},
			{
				ResourceName:            "google_developer_connect_connection.my-connection",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"connection_id", "location", "terraform_labels"},
			},
		},
	})
}

func testAccDeveloperConnectConnection_Gitlab(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_developer_connect_connection" "my-connection" {
  location = "us-central1"
  connection_id = "tf-test-tf-test-connection%{random_suffix}"

  gitlab_config {
    webhook_secret_secret_version = "projects/devconnect-terraform-creds/secrets/gitlab-webhook/versions/latest"

    read_authorizer_credential {
      user_token_secret_version = "projects/devconnect-terraform-creds/secrets/gitlab-read-cred/versions/latest"
    }

    authorizer_credential {
      user_token_secret_version = "projects/devconnect-terraform-creds/secrets/gitlab-auth-cred/versions/latest"
    }
  }
}
`, context)
}

func testAccDeveloperConnectConnection_GitlabUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_developer_connect_connection" "my-connection" {
  location = "us-central1"
  connection_id = "tf-test-tf-test-connection%{random_suffix}"
  annotations = {}
  labels = {}
  
  crypto_key_config {
    key_reference = "projects/devconnect-terraform-creds/locations/us-central1/keyRings/tf-keyring/cryptoKeys/tf-crypto-key"
  }

  gitlab_config {
    webhook_secret_secret_version = "projects/devconnect-terraform-creds/secrets/gitlab-webhook/versions/latest"

    read_authorizer_credential {
      user_token_secret_version = "projects/devconnect-terraform-creds/secrets/gitlab-read-cred-update/versions/latest"
    }

    authorizer_credential {
      user_token_secret_version = "projects/devconnect-terraform-creds/secrets/gitlab-auth-cred-update/versions/latest"
    }
  }
}
`, context)
}

func TestAccDeveloperConnectConnection_GlePrivConnection(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDeveloperConnectConnection_GlePrivConnection(context),
			},
			{
				ResourceName:            "google_developer_connect_connection.my-connection",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"connection_id", "location", "terraform_labels"},
			},
		},
	})
}

func testAccDeveloperConnectConnection_GlePrivConnection(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_developer_connect_connection" "my-connection" {
  location = "us-central1"
  connection_id = "tf-test-tf-test-connection%{random_suffix}"
  annotations = {}
  labels = {}

  gitlab_enterprise_config {
    host_uri = "https://gle-us.gle-us-private.com"

    webhook_secret_secret_version = "projects/devconnect-terraform-creds/secrets/gitlab-enterprise-webhook/versions/latest"

    read_authorizer_credential {
      user_token_secret_version = "projects/devconnect-terraform-creds/secrets/gle-private-read-token/versions/latest"
    }

    authorizer_credential {
      user_token_secret_version = "projects/devconnect-terraform-creds/secrets/gle-private-api/versions/latest"
    }

    ssl_ca_certificate = "-----BEGIN CERTIFICATE-----\nMIIFbjCCA1agAwIBAgIUH+nsWsqagMW9Ld8E9J71yPLPpD8wDQYJKoZIhvcNAQEL\nBQAwJDEiMCAGA1UEAwwZZ2xlLXVzLmdsZS11cy1wcml2YXRlLmNvbTAeFw0yNDEw\nMzExNjQzMjBaFw0zNDEwMjkxNjQzMjBaMCQxIjAgBgNVBAMMGWdsZS11cy5nbGUt\ndXMtcHJpdmF0ZS5jb20wggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQDL\n+dUU8MHo+Eskx4SSnI1thRIiljgsyJSzSOplaD4lqahFnrG0cB0ovKpyRL4A+0wM\nzVW7W1Pfi8DiEOxxfNo7pEj+0zrzJHHqnzW9kApIlRmO1TBBJ7i9HaVamJ1Od01b\n2WI/pnKFEvNfLQSDQBulhkHZ2McyauDhb1DxefKnVX8ac6qhxtc4IzrezIQuJ18N\nDPtNLUDD4rtU4mIX4lx1yBIplrgypAo1HDbJOwW8OR76MtjAY7ek3K6UCyduQtwy\nmfZ23b3Eg69W10HVMVTy6m5NaGKi/TWy2MJ71hKUQ1+tWIPH5SL7FzYPKL4XXw5W\np61LhIiBAd2tgD41b2cQxhUbVifc1qHtnwNz/tE77M9ySH37rEUIlExzr3D3JV+f\nXjXEXUr9as8GRnS5zhD/opKe7wKbwpYMHhylK1h6XH/sBO7dBajf5xCvpZZBDzrK\nkpTqwHspT7p40WF9d8odjEk/xZKn5LdcDG2I+4U7SVS5e8ud41HUQxJwQx56lKfh\n2WB+zs7nSyMfspTj4doY1OADEC1VQCyGrwlbclKTKmUWrgwQdm38KxDzW5Juyjmm\nzvfsWIlSMdnes0qVVo38N3Jz8/MHCLD79R2veWgA2fbqS3+4h2dRkR7htjaVlJMJ\nt7SwFiG39ic3OZpo+wTkaHlG4CBnbFDueUsOW2wEpQIDAQABo4GXMIGUMB0GA1Ud\nDgQWBBTExgzH2gz9+rJHvlTFPO0AvG88azAfBgNVHSMEGDAWgBTExgzH2gz9+rJH\nvlTFPO0AvG88azAPBgNVHRMBAf8EBTADAQH/MEEGA1UdEQQ6MDiCGWdsZS11cy5n\nbGUtdXMtcHJpdmF0ZS5jb22CGyouZ2xlLXVzLmdsZS11cy1wcml2YXRlLmNvbTAN\nBgkqhkiG9w0BAQsFAAOCAgEAjkd1ZNoekoWrmozD+Ta1OM0zWhv04eqhP8aYzhbd\nXRS+GyF6ifMwfWg9HogkH22ZPT5GszaL5DacSyOUqZgJ905Q6g1EFPnaKmFVHHeC\nzZAhg5oedAzcakZpYwZDSiLuPgsQfwgRnqWIYR8JcIM5bKRZNGyOg8eZ8cKu23A2\nPavL4B3Ra1l93KllKm21rigIhLPIPLoEyxEg9c9oTJF92r0+aRdf2Ln853260Fqf\ncEUWoXhqMGvDv/YEbqDjGQ/Kh7ZWdlIWhcKFOA0gluF7oExjt/MgSitukgg3aaic\n/eXXOrZDNYH7Ve610NUuNlhub1M47Tp7EgjUJVWlsKK84T8ZcZq7Hn4BzioUr95d\nHao6u19HWA/ISM8bwzHaYxscFI4u6phEL0HJzLf4EysEmS0rAnLxyol0apNx6znR\nhXsqxnSexKhXoLqnK1Vuhcg8DsvobXHqg68EGZ7BZ3ycPYaHSWU8Xh3l1gtYkcQ6\nzxXsKIijlpVKuYJvGA3EOMoZu6+2MYF8Tgp3N4sKMvPhqBhsmgxOYF5OkAbGlsUP\nyCYWFDBFHmbhvUu5JpbKuID2CPkBi16EetemvMQ9PGlLq/0fO/BBNkn6TYn9Kvg8\nAyvuONz54uFEAIKPCcZIosa3ml+5/pt+tBhtVzHA6vMxn18IYaNpuTwSxi/+M10K\nRjw=\n-----END CERTIFICATE-----\n"

       
    service_directory_config {
      service = "projects/devconnect-terraform-creds/locations/us-central1/namespaces/my-namespace/services/terraform-gle"
    }
  }

}
`, context)
}

func TestAccDeveloperConnectConnection_developerConnectConnectionGitlabEnterpriseUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDeveloperConnectConnection_GitlabEnterprise(context),
			},
			{
				ResourceName:            "google_developer_connect_connection.my-connection",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"connection_id", "location", "terraform_labels"},
			},
			{
				Config: testAccDeveloperConnectConnection_GitlabEnterpriseUpdate(context),
			},
			{
				ResourceName:            "google_developer_connect_connection.my-connection",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"connection_id", "location", "terraform_labels"},
			},
		},
	})
}

func testAccDeveloperConnectConnection_GitlabEnterprise(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_developer_connect_connection" "my-connection" {
  location = "us-central1"
  connection_id = "tf-test-tf-test-connection%{random_suffix}"

  gitlab_enterprise_config {
    host_uri = "https://gle-us-central1.gcb-test.com"

    webhook_secret_secret_version = "projects/devconnect-terraform-creds/secrets/gitlab-enterprise-webhook/versions/latest"

    read_authorizer_credential {
      user_token_secret_version = "projects/devconnect-terraform-creds/secrets/gitlab-enterprise-read-cred/versions/latest"
    }

    authorizer_credential {
      user_token_secret_version = "projects/devconnect-terraform-creds/secrets/gitlab-enterprise-auth-cred/versions/latest"
    }
  }
}
`, context)
}

func testAccDeveloperConnectConnection_GitlabEnterpriseUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_developer_connect_connection" "my-connection" {
  location = "us-central1"
  connection_id = "tf-test-tf-test-connection%{random_suffix}"
  annotations = {}
  labels = {}
  
  crypto_key_config {
    key_reference = "projects/devconnect-terraform-creds/locations/us-central1/keyRings/tf-keyring/cryptoKeys/tf-crypto-key"
  }

  gitlab_enterprise_config {
    host_uri = "https://gle-old.gcb-test.com"

    webhook_secret_secret_version = "projects/devconnect-terraform-creds/secrets/gitlab-enterprise-webhook/versions/latest"

    read_authorizer_credential {
      user_token_secret_version = "projects/devconnect-terraform-creds/secrets/gitlab-enterprise-read-cred-update/versions/latest"
    }

    authorizer_credential {
      user_token_secret_version = "projects/devconnect-terraform-creds/secrets/gitlab-enterprise-auth-cred-update/versions/latest"
    }
  }
}
`, context)
}
