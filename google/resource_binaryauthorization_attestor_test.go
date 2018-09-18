package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccBinaryAuthorizationAttestor_basic(t *testing.T) {
	t.Parallel()

	name := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBinaryAuthorizationAttestorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBinaryAuthorizationAttestorBasic(name),
			},
			{
				ResourceName:      "google_binary_authorization_attestor.attestor",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBinaryAuthorizationAttestor_full(t *testing.T) {
	t.Parallel()

	name := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBinaryAuthorizationAttestorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBinaryAuthorizationAttestorFull(name),
			},
			{
				ResourceName:      "google_binary_authorization_attestor.attestor",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBinaryAuthorizationAttestor_update(t *testing.T) {
	t.Parallel()

	name := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBinaryAuthorizationAttestorDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccBinaryAuthorizationAttestorBasic(name),
			},
			{
				ResourceName:      "google_binary_authorization_attestor.attestor",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBinaryAuthorizationAttestorFull(name),
			},
			{
				ResourceName:      "google_binary_authorization_attestor.attestor",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBinaryAuthorizationAttestorBasic(name),
			},
			{
				ResourceName:      "google_binary_authorization_attestor.attestor",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckBinaryAuthorizationAttestorDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_binary_authorization_attestor" {
			continue
		}

		project, err := getTestProject(rs.Primary, config)
		if err != nil {
			return err
		}

		name := rs.Primary.Attributes["name"]

		url := fmt.Sprintf("https://binaryauthorization.googleapis.com/v1beta1/projects/%s/attestors/%s", project, name)
		_, err = sendRequest(config, "GET", url, nil)

		if err == nil {
			return fmt.Errorf("Error, attestor %s still exists", name)
		}
	}

	return nil
}

func testAccBinaryAuthorizationAttestorBasic(name string) string {
	return fmt.Sprintf(`
resource "google_container_analysis_note" "note" {
  name = "tf-test-%s"
  attestation_authority {
    hint {
      human_readable_name = "My Attestor"
    }
  }
}

resource "google_binary_authorization_attestor" "attestor" {
  name = "tf-test-%s"
  attestation_authority_note {
    note_reference = "${google_container_analysis_note.note.name}"
  }
}
`, name, name)
}

func testAccBinaryAuthorizationAttestorFull(name string) string {
	return fmt.Sprintf(`
resource "google_container_analysis_note" "note" {
  name = "tf-test-%s"
  attestation_authority {
    hint {
      human_readable_name = "My Attestor"
    }
  }
}

resource "google_binary_authorization_attestor" "attestor" {
  name = "tf-test-%s"
  description = "my description"
  attestation_authority_note {
    note_reference = "${google_container_analysis_note.note.name}"
    public_keys {
      ascii_armored_pgp_public_key = <<EOF
%s
EOF
      comment = "this key has a comment"
    }
  }
}
`, name, name, armoredPubKey)
}

// Generated key using instructions from
// https://cloud.google.com/binary-authorization/docs/creating-attestors#generate_pgp_key_pairs.
// This key has no real meaning and was generated in order to have a valid key
// for testing.
const armoredPubKey = `mQENBFtP0doBCADF+joTiXWKVuP8kJt3fgpBSjT9h8ezMfKA4aXZctYLx5wslWQl
bB7Iu2ezkECNzoEeU7WxUe8a61pMCh9cisS9H5mB2K2uM4Jnf8tgFeXn3akJDVo0
oR1IC+Dp9mXbRSK3MAvKkOwWlG99sx3uEdvmeBRHBOO+grchLx24EThXFOyP9Fk6
V39j6xMjw4aggLD15B4V0v9JqBDdJiIYFzszZDL6pJwZrzcP0z8JO4rTZd+f64bD
Mpj52j/pQfA8lZHOaAgb1OrthLdMrBAjoDjArV4Ek7vSbrcgYWcI6BhsQrFoxKdX
83TZKai55ZCfCLIskwUIzA1NLVwyzCS+fSN/ABEBAAG0KCJUZXN0IEF0dGVzdG9y
IiA8ZGFuYWhvZmZtYW5AZ29vZ2xlLmNvbT6JAU4EEwEIADgWIQRfWkqHt6hpTA1L
uY060eeM4dc66AUCW0/R2gIbLwULCQgHAgYVCgkICwIEFgIDAQIeAQIXgAAKCRA6
0eeM4dc66HdpCAC4ot3b0OyxPb0Ip+WT2U0PbpTBPJklesuwpIrM4Lh0N+1nVRLC
51WSmVbM8BiAFhLbN9LpdHhds1kUrHF7+wWAjdR8sqAj9otc6HGRM/3qfa2qgh+U
WTEk/3us/rYSi7T7TkMuutRMIa1IkR13uKiW56csEMnbOQpn9rDqwIr5R8nlZP5h
MAU9vdm1DIv567meMqTaVZgR3w7bck2P49AO8lO5ERFpVkErtu/98y+rUy9d789l
+OPuS1NGnxI1YKsNaWJF4uJVuvQuZ1twrhCbGNtVorO2U12+cEq+YtUxj7kmdOC1
qoIRW6y0+UlAc+MbqfL0ziHDOAmcqz1GnROg
=6Bvm`
