// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package integrations_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccIntegrationsAuthConfig_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIntegrationsAuthConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationsAuthConfig_full(context),
			},
			{
				ResourceName:            "google_integrations_auth_config.update_example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"client_certificate", "location"},
			},
			{
				Config: testAccIntegrationsAuthConfig_update(context),
			},
			{
				ResourceName:            "google_integrations_auth_config.update_example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"client_certificate", "location"},
			},
		},
	})
}

func testAccIntegrationsAuthConfig_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_integrations_client" "client" {
	location = "southamerica-west1"
	provision_gmek = true
}

resource "google_integrations_auth_config" "update_example" {
    location = "southamerica-west1"
    display_name = "tf-test-test-authconfig%{random_suffix}"
    description = "Test auth config created via terraform"
    visibility = "CLIENT_VISIBLE"
    expiry_notification_duration = ["3.500s"]
    override_valid_time = "2014-10-02T15:01:23Z"
    decrypted_credential {
        credential_type = "USERNAME_AND_PASSWORD"
        username_and_password {
            username = "test-username"
            password = "test-password"
        }
    }
    depends_on = [google_integrations_client.client]
}
`, context)
}

func testAccIntegrationsAuthConfig_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_integrations_client" "client" {
	location = "southamerica-west1"
	provision_gmek = true
}

resource "google_integrations_auth_config" "update_example" {
    location = "southamerica-west1"
    display_name = "tf-test-test-authconfig-update%{random_suffix}"
    description = "Test auth config updated via terraform"
    visibility = "CLIENT_VISIBLE"
    expiry_notification_duration = ["4s"]
    override_valid_time = "2014-10-10T15:01:23Z"
    decrypted_credential {
        credential_type = "CLIENT_CERTIFICATE_ONLY"
    }
    client_certificate {
        ssl_certificate = <<EOT
-----BEGIN CERTIFICATE-----
MIICTTCCAbagAwIBAgIJAPT0tSKNxan/MA0GCSqGSIb3DQEBCwUAMCoxFzAVBgNV
BAoTDkdvb2dsZSBURVNUSU5HMQ8wDQYDVQQDEwZ0ZXN0Q0EwHhcNMTUwMTAxMDAw
MDAwWhcNMjUwMTAxMDAwMDAwWjAuMRcwFQYDVQQKEw5Hb29nbGUgVEVTVElORzET
MBEGA1UEAwwKam9lQGJhbmFuYTCBnzANBgkqhkiG9w0BAQEFAAOBjQAwgYkCgYEA
vDYFgMgxi5W488d9J7UpCInl0NXmZQpJDEHE4hvkaRlH7pnC71H0DLt0/3zATRP1
JzY2+eqBmbGl4/sgZKYv8UrLnNyQNUTsNx1iZAfPUflf5FwgVsai8BM0pUciq1NB
xD429VFcrGZNucvFLh72RuRFIKH8WUpiK/iZNFkWhZ0CAwEAAaN3MHUwDgYDVR0P
AQH/BAQDAgWgMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAMBgNVHRMB
Af8EAjAAMBkGA1UdDgQSBBCVgnFBCWgL/iwCqnGrhTPQMBsGA1UdIwQUMBKAEKey
Um2o4k2WiEVA0ldQvNYwDQYJKoZIhvcNAQELBQADgYEAYK986R4E3L1v+Q6esBtW
JrUwA9UmJRSQr0N5w3o9XzarU37/bkjOP0Fw0k/A6Vv1n3vlciYfBFaBIam1qRHr
5dMsYf4CZS6w50r7hyzqyrwDoyNxkLnd2PdcHT/sym1QmflsjEs7pejtnohO6N2H
wQW6M0H7Zt8claGRla4fKkg=
-----END CERTIFICATE-----
EOT
		encrypted_private_key = <<EOT
-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCA/Oj2HXqs5fTk
j/8DrlOQtLG3K9RMsYHvnwICLxkGqVcTfut58hDFLbQM8C3C0ENAKitNJplCJmYG
8VpgZzgq8VxaGnlP/sXUFLMGksd5sATn0sY3SkPndTKk/dqqA4MIh/dYfh19ynEN
hB9Ll/h54Yic2je2Qaxe/uMMu8RODTz3oCn7FcoYpPvfygfU0ntn4IcqH/hts5DG
s+3otJk4entRZglQDxR+sWOsbLtJIQZDP8rH3jDVdl5l3wspgtMTY8b5T5+pLm0p
/OzCmxT0dq/O6BhpxI1xf/zcdRZeWk5DTJxTi5AgPquTlAG/B6A3HkqBJ14hT/Rk
iv7Ma3DLAgMBAAECggEABATkf9VfpiAT9zYdouk50bBpckvymQTyQLD8SlBaX+KY
kgv/pHSXK4Pm4iensrQerFLgfqPA3U+FiqjW5Mv7c1VRK6HJbuVkpdzoXLI9IQsL
vsBY7//9Ajk5P7NokjdB6JPdU/2dHROuQVa59cxPtzpHo0htnPlDOKXfFZZuoZ17
Nr8WQHrHy8P8ABM1tLOzvU9Nlh7TcjQvev+HxkLek4qzYyJ/Ac7XOjg/XKUm1tZk
O3BHr8YLabwyjO7l1t+2b14rUTL/8pfUZnAkEi3FAlPxm3ilftmX65zliC9G4ghk
dr5PByT3DqnuIIglua9bISv1H34ogecd+9a6EU7RxQKBgQC2RPKLounXZo8vYiU4
sFTEvjbs+u9Ypk4OrNLnb8KdacLBUaJGnf++xbBoKpwFCBJfy//fvuQfusYF9Gyn
GxL43tw94C/H5upQYnDsmnQak6TbOu3mA24OGK7Rcq6NEHgeCY4HomutnSiPTZJq
8jlpqgqh1itETe5avgkMNq3zBwKBgQC1KlztGzvbB+rUDc6Kfvk5pUbCSFKMMMa2
NWNXeD6i2iA56zEYSbTjKQ3u9pjUV8LNqAdUFxmbdPxZjheNK2dEm68SVRXPKOeB
EmQT+t/EyW9LqBEA2oZt3h2hXtK8ppJjQm4XUCDs1NphP87eNzx5FLzJWjG8VqDq
jOvApNqPHQKBgDQqlZSbgvvwUYjJOUf5R7mri0LWKwyfRHX0xsQQe43cCC6WM7Cs
Zdbu86dMkqzp+4BJfalHFDl0llp782D8Ybiy6CwZbvNyxptNIW7GYfZ9TVCllBMh
5izIqbgub4DWNtq591l+Bf2BnmstU3uiagYw8awSBP4eo9p6y1IgkDafAoGBAJbi
lIiqEP0IqA06/pWc0Qew3rD7OT0ndqjU6Es2i7xovURf3QDkinJThBZNbdYUzdsp
IgloP9yY33/a90SNLLIYlARJtyNVZxK59X4qiOpF9prlfFvgpOumfbkj15JljTB8
aGKkSvfVA5jRYwLysDwMCHwO0bOR1u3itos5AgsFAoGAKEGms1kuQ5/HyFgSmg9G
wBUzu+5Y08/A37rvyXsR6GjmlZJvULEopJNUNCOOpITNQikXK63sIFry7/59eGv5
UwKadZbfwbVF5ipu59UxfVE3lipf/mYePDqMkHVWv/8p+OnnJt9uKnyW8VSOu5uk
82QF30zbIWDTUjrcugVAs+E=
-----END PRIVATE KEY-----     
EOT
		passphrase = ""
    }
    depends_on = [google_integrations_client.client]
}
`, context)
}
