// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccBinaryAuthorizationAttestorIamBindingGenerated(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": RandString(t, 10),
		"role":          "roles/viewer",
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBinaryAuthorizationAttestorIamBinding_basicGenerated(context),
			},
			{
				ResourceName:      "google_binary_authorization_attestor_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/attestors/%s roles/viewer", acctest.GetTestProjectFromEnv(), fmt.Sprintf("tf-test-test-attestor%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test Iam Binding update
				Config: testAccBinaryAuthorizationAttestorIamBinding_updateGenerated(context),
			},
			{
				ResourceName:      "google_binary_authorization_attestor_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/attestors/%s roles/viewer", acctest.GetTestProjectFromEnv(), fmt.Sprintf("tf-test-test-attestor%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBinaryAuthorizationAttestorIamMemberGenerated(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": RandString(t, 10),
		"role":          "roles/viewer",
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test Iam Member creation (no update for member, no need to test)
				Config: testAccBinaryAuthorizationAttestorIamMember_basicGenerated(context),
			},
			{
				ResourceName:      "google_binary_authorization_attestor_iam_member.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/attestors/%s roles/viewer user:admin@hashicorptest.com", acctest.GetTestProjectFromEnv(), fmt.Sprintf("tf-test-test-attestor%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBinaryAuthorizationAttestorIamPolicyGenerated(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": RandString(t, 10),
		"role":          "roles/viewer",
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBinaryAuthorizationAttestorIamPolicy_basicGenerated(context),
				Check:  resource.TestCheckResourceAttrSet("data.google_binary_authorization_attestor_iam_policy.foo", "policy_data"),
			},
			{
				ResourceName:      "google_binary_authorization_attestor_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/attestors/%s", acctest.GetTestProjectFromEnv(), fmt.Sprintf("tf-test-test-attestor%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBinaryAuthorizationAttestorIamPolicy_emptyBinding(context),
			},
			{
				ResourceName:      "google_binary_authorization_attestor_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/attestors/%s", acctest.GetTestProjectFromEnv(), fmt.Sprintf("tf-test-test-attestor%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccBinaryAuthorizationAttestorIamMember_basicGenerated(context map[string]interface{}) string {
	return Nprintf(`
resource "google_binary_authorization_attestor" "attestor" {
  name = "tf-test-test-attestor%{random_suffix}"
  attestation_authority_note {
    note_reference = google_container_analysis_note.note.name
    public_keys {
      ascii_armored_pgp_public_key = <<EOF
mQENBFtP0doBCADF+joTiXWKVuP8kJt3fgpBSjT9h8ezMfKA4aXZctYLx5wslWQl
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
=6Bvm
EOF

    }
  }
}

resource "google_container_analysis_note" "note" {
  name = "tf-test-test-attestor-note%{random_suffix}"
  attestation_authority {
    hint {
      human_readable_name = "Attestor Note"
    }
  }
}

resource "google_binary_authorization_attestor_iam_member" "foo" {
  project = google_binary_authorization_attestor.attestor.project
  attestor = google_binary_authorization_attestor.attestor.name
  role = "%{role}"
  member = "user:admin@hashicorptest.com"
}
`, context)
}

func testAccBinaryAuthorizationAttestorIamPolicy_basicGenerated(context map[string]interface{}) string {
	return Nprintf(`
resource "google_binary_authorization_attestor" "attestor" {
  name = "tf-test-test-attestor%{random_suffix}"
  attestation_authority_note {
    note_reference = google_container_analysis_note.note.name
    public_keys {
      ascii_armored_pgp_public_key = <<EOF
mQENBFtP0doBCADF+joTiXWKVuP8kJt3fgpBSjT9h8ezMfKA4aXZctYLx5wslWQl
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
=6Bvm
EOF

    }
  }
}

resource "google_container_analysis_note" "note" {
  name = "tf-test-test-attestor-note%{random_suffix}"
  attestation_authority {
    hint {
      human_readable_name = "Attestor Note"
    }
  }
}

data "google_iam_policy" "foo" {
  binding {
    role = "%{role}"
    members = ["user:admin@hashicorptest.com"]
  }
}

resource "google_binary_authorization_attestor_iam_policy" "foo" {
  project = google_binary_authorization_attestor.attestor.project
  attestor = google_binary_authorization_attestor.attestor.name
  policy_data = data.google_iam_policy.foo.policy_data
}

data "google_binary_authorization_attestor_iam_policy" "foo" {
  project = google_binary_authorization_attestor.attestor.project
  attestor = google_binary_authorization_attestor.attestor.name
  depends_on = [
    google_binary_authorization_attestor_iam_policy.foo
  ]
}
`, context)
}

func testAccBinaryAuthorizationAttestorIamPolicy_emptyBinding(context map[string]interface{}) string {
	return Nprintf(`
resource "google_binary_authorization_attestor" "attestor" {
  name = "tf-test-test-attestor%{random_suffix}"
  attestation_authority_note {
    note_reference = google_container_analysis_note.note.name
    public_keys {
      ascii_armored_pgp_public_key = <<EOF
mQENBFtP0doBCADF+joTiXWKVuP8kJt3fgpBSjT9h8ezMfKA4aXZctYLx5wslWQl
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
=6Bvm
EOF

    }
  }
}

resource "google_container_analysis_note" "note" {
  name = "tf-test-test-attestor-note%{random_suffix}"
  attestation_authority {
    hint {
      human_readable_name = "Attestor Note"
    }
  }
}

data "google_iam_policy" "foo" {
}

resource "google_binary_authorization_attestor_iam_policy" "foo" {
  project = google_binary_authorization_attestor.attestor.project
  attestor = google_binary_authorization_attestor.attestor.name
  policy_data = data.google_iam_policy.foo.policy_data
}
`, context)
}

func testAccBinaryAuthorizationAttestorIamBinding_basicGenerated(context map[string]interface{}) string {
	return Nprintf(`
resource "google_binary_authorization_attestor" "attestor" {
  name = "tf-test-test-attestor%{random_suffix}"
  attestation_authority_note {
    note_reference = google_container_analysis_note.note.name
    public_keys {
      ascii_armored_pgp_public_key = <<EOF
mQENBFtP0doBCADF+joTiXWKVuP8kJt3fgpBSjT9h8ezMfKA4aXZctYLx5wslWQl
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
=6Bvm
EOF

    }
  }
}

resource "google_container_analysis_note" "note" {
  name = "tf-test-test-attestor-note%{random_suffix}"
  attestation_authority {
    hint {
      human_readable_name = "Attestor Note"
    }
  }
}

resource "google_binary_authorization_attestor_iam_binding" "foo" {
  project = google_binary_authorization_attestor.attestor.project
  attestor = google_binary_authorization_attestor.attestor.name
  role = "%{role}"
  members = ["user:admin@hashicorptest.com"]
}
`, context)
}

func testAccBinaryAuthorizationAttestorIamBinding_updateGenerated(context map[string]interface{}) string {
	return Nprintf(`
resource "google_binary_authorization_attestor" "attestor" {
  name = "tf-test-test-attestor%{random_suffix}"
  attestation_authority_note {
    note_reference = google_container_analysis_note.note.name
    public_keys {
      ascii_armored_pgp_public_key = <<EOF
mQENBFtP0doBCADF+joTiXWKVuP8kJt3fgpBSjT9h8ezMfKA4aXZctYLx5wslWQl
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
=6Bvm
EOF

    }
  }
}

resource "google_container_analysis_note" "note" {
  name = "tf-test-test-attestor-note%{random_suffix}"
  attestation_authority {
    hint {
      human_readable_name = "Attestor Note"
    }
  }
}

resource "google_binary_authorization_attestor_iam_binding" "foo" {
  project = google_binary_authorization_attestor.attestor.project
  attestor = google_binary_authorization_attestor.attestor.name
  role = "%{role}"
  members = ["user:admin@hashicorptest.com", "user:gterraformtest1@gmail.com"]
}
`, context)
}
