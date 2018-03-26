package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// Test that a service account key can be created and destroyed
func TestAccServiceAccountKey_basic(t *testing.T) {
	t.Parallel()

	resourceName := "google_service_account_key.acceptance"
	accountID := "a" + acctest.RandString(10)
	displayName := "Terraform Test"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccServiceAccountKey(accountID, displayName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleServiceAccountKeyExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "public_key"),
					resource.TestCheckResourceAttrSet(resourceName, "valid_after"),
					resource.TestCheckResourceAttrSet(resourceName, "valid_before"),
					resource.TestCheckResourceAttrSet(resourceName, "private_key"),
				),
			},
		},
	})
}

func TestAccServiceAccountKey_fromEmail(t *testing.T) {
	t.Parallel()

	resourceName := "google_service_account_key.acceptance"
	accountID := "a" + acctest.RandString(10)
	displayName := "Terraform Test"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccServiceAccountKey_fromEmail(accountID, displayName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleServiceAccountKeyExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "public_key"),
					resource.TestCheckResourceAttrSet(resourceName, "valid_after"),
					resource.TestCheckResourceAttrSet(resourceName, "valid_before"),
					resource.TestCheckResourceAttrSet(resourceName, "private_key"),
				),
			},
		},
	})
}

func TestAccServiceAccountKey_pgp(t *testing.T) {
	t.Parallel()
	resourceName := "google_service_account_key.acceptance"
	accountID := "a" + acctest.RandString(10)
	displayName := "Terraform Test"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccServiceAccountKey_pgp(accountID, displayName, testKeyPairPubKey1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleServiceAccountKeyExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "public_key"),
					resource.TestCheckResourceAttrSet(resourceName, "private_key_encrypted"),
					resource.TestCheckResourceAttrSet(resourceName, "private_key_fingerprint"),
				),
			},
		},
	})
}

func testAccCheckGoogleServiceAccountKeyExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("Not found: %s", r)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		config := testAccProvider.Meta().(*Config)

		_, err := config.clientIAM.Projects.ServiceAccounts.Keys.Get(rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccServiceAccountKey(account, name string) string {
	return fmt.Sprintf(`
resource "google_service_account" "acceptance" {
	account_id = "%s"
	display_name = "%s"
}

resource "google_service_account_key" "acceptance" {
	service_account_id = "${google_service_account.acceptance.name}"
	public_key_type = "TYPE_X509_PEM_FILE"
}
`, account, name)
}

func testAccServiceAccountKey_fromEmail(account, name string) string {
	return fmt.Sprintf(`
resource "google_service_account" "acceptance" {
	account_id = "%s"
	display_name = "%s"
}

resource "google_service_account_key" "acceptance" {
	service_account_id = "${google_service_account.acceptance.email}"
	public_key_type = "TYPE_X509_PEM_FILE"
}
`, account, name)
}

func testAccServiceAccountKey_pgp(account, name string, key string) string {
	return fmt.Sprintf(`
resource "google_service_account" "acceptance" {
	account_id = "%s"
	display_name = "%s"
}

resource "google_service_account_key" "acceptance" {
	service_account_id = "${google_service_account.acceptance.name}"
	public_key_type = "TYPE_X509_PEM_FILE"
	pgp_key = <<EOF
%s
EOF
}
`, account, name, key)
}

const testKeyPairPubKey1 = `mQENBFXbjPUBCADjNjCUQwfxKL+RR2GA6pv/1K+zJZ8UWIF9S0lk7cVIEfJiprzzwiMwBS5cD0da
rGin1FHvIWOZxujA7oW0O2TUuatqI3aAYDTfRYurh6iKLC+VS+F7H+/mhfFvKmgr0Y5kDCF1j0T/
063QZ84IRGucR/X43IY7kAtmxGXH0dYOCzOe5UBX1fTn3mXGe2ImCDWBH7gOViynXmb6XNvXkP0f
sF5St9jhO7mbZU9EFkv9O3t3EaURfHopsCVDOlCkFCw5ArY+DUORHRzoMX0PnkyQb5OzibkChzpg
8hQssKeVGpuskTdz5Q7PtdW71jXd4fFVzoNH8fYwRpziD2xNvi6HABEBAAG0EFZhdWx0IFRlc3Qg
S2V5IDGJATgEEwECACIFAlXbjPUCGy8GCwkIBwMCBhUIAgkKCwQWAgMBAh4BAheAAAoJEOfLr44B
HbeTo+sH/i7bapIgPnZsJ81hmxPj4W12uvunksGJiC7d4hIHsG7kmJRTJfjECi+AuTGeDwBy84TD
cRaOB6e79fj65Fg6HgSahDUtKJbGxj/lWzmaBuTzlN3CEe8cMwIPqPT2kajJVdOyrvkyuFOdPFOE
A7bdCH0MqgIdM2SdF8t40k/ATfuD2K1ZmumJ508I3gF39jgTnPzD4C8quswrMQ3bzfvKC3klXRlB
C0yoArn+0QA3cf2B9T4zJ2qnvgotVbeK/b1OJRNj6Poeo+SsWNc/A5mw7lGScnDgL3yfwCm1gQXa
QKfOt5x+7GqhWDw10q+bJpJlI10FfzAnhMF9etSqSeURBRW5AQ0EVduM9QEIAL53hJ5bZJ7oEDCn
aY+SCzt9QsAfnFTAnZJQrvkvusJzrTQ088eUQmAjvxkfRqnv981fFwGnh2+I1Ktm698UAZS9Jt8y
jak9wWUICKQO5QUt5k8cHwldQXNXVXFa+TpQWQR5yW1a9okjh5o/3d4cBt1yZPUJJyLKY43Wvptb
6EuEsScO2DnRkh5wSMDQ7dTooddJCmaq3LTjOleRFQbu9ij386Do6jzK69mJU56TfdcydkxkWF5N
ZLGnED3lq+hQNbe+8UI5tD2oP/3r5tXKgMy1R/XPvR/zbfwvx4FAKFOP01awLq4P3d/2xOkMu4Lu
9p315E87DOleYwxk+FoTqXEAEQEAAYkCPgQYAQIACQUCVduM9QIbLgEpCRDny6+OAR23k8BdIAQZ
AQIABgUCVduM9QAKCRAID0JGyHtSGmqYB/4m4rJbbWa7dBJ8VqRU7ZKnNRDR9CVhEGipBmpDGRYu
lEimOPzLUX/ZXZmTZzgemeXLBaJJlWnopVUWuAsyjQuZAfdd8nHkGRHG0/DGum0l4sKTta3OPGHN
C1z1dAcQ1RCr9bTD3PxjLBczdGqhzw71trkQRBRdtPiUchltPMIyjUHqVJ0xmg0hPqFic0fICsr0
YwKoz3h9+QEcZHvsjSZjgydKvfLYcm+4DDMCCqcHuJrbXJKUWmJcXR0y/+HQONGrGJ5xWdO+6eJi
oPn2jVMnXCm4EKc7fcLFrz/LKmJ8seXhxjM3EdFtylBGCrx3xdK0f+JDNQaC/rhUb5V2XuX6VwoH
/AtY+XsKVYRfNIupLOUcf/srsm3IXT4SXWVomOc9hjGQiJ3rraIbADsc+6bCAr4XNZS7moViAAcI
PXFv3m3WfUlnG/om78UjQqyVACRZqqAGmuPq+TSkRUCpt9h+A39LQWkojHqyob3cyLgy6z9Q557O
9uK3lQozbw2gH9zC0RqnePl+rsWIUU/ga16fH6pWc1uJiEBt8UZGypQ/E56/343epmYAe0a87sHx
8iDV+dNtDVKfPRENiLOOc19MmS+phmUyrbHqI91c0pmysYcJZCD3a502X1gpjFbPZcRtiTmGnUKd
OIu60YPNE4+h7u2CfYyFPu3AlUaGNMBlvy6PEpU=`
