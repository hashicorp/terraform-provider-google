package google

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBillingBudget_billingBudgetCurrencycode(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"billing_acct":  getTestBillingAccountFromEnv(t),
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBillingBudgetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBillingBudget_billingBudgetCurrencycode(context),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func testAccBillingBudget_billingBudgetCurrencycode(context map[string]interface{}) string {
	return Nprintf(`
data "google_billing_account" "account" {
  billing_account = "%{billing_acct}"
}

data "google_project" "project" {
}

resource "google_billing_budget" "budget" {
  billing_account = data.google_billing_account.account.id
  display_name    = "Example Billing Budget%{random_suffix}"

  budget_filter {
    projects = ["projects/${data.google_project.project.number}"]
  }

  amount {
    specified_amount {
      units         = "100000"
    }
  }

  threshold_rules {
    threshold_percent = 1.0
  }
  threshold_rules {
    threshold_percent = 1.0
    spend_basis       = "FORECASTED_SPEND"
  }
}
`, context)
}

func TestAccBillingBudget_billingBudgetUpdateRemoveFilter(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"billing_acct":  getTestBillingAccountFromEnv(t),
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBillingBudgetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBillingBudget_billingBudgetUpdateRemoveFilterStart(context),
			},
			{
				ResourceName:      "google_billing_budget.budget",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBillingBudget_billingBudgetUpdateRemoveFilterEnd(context),
			},
			{
				ResourceName:      "google_billing_budget.budget",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccBillingBudget_billingBudgetUpdateRemoveFilterStart(context map[string]interface{}) string {
	return Nprintf(`
data "google_billing_account" "account" {
  billing_account = "%{billing_acct}"
}

data "google_project" "project" {
}

resource "google_billing_budget" "budget" {
  billing_account = data.google_billing_account.account.id
  display_name = "Example Billing Budget%{random_suffix}"

  budget_filter {
    projects = ["projects/${data.google_project.project.number}"]
  }

  amount {
    specified_amount {
      currency_code = "USD"
      units = "100000"
    }
  }

  threshold_rules {
    threshold_percent = 0.5
  }
  threshold_rules {
    threshold_percent = 0.9
    spend_basis = "FORECASTED_SPEND"
  }
}
`, context)
}

func testAccBillingBudget_billingBudgetUpdateRemoveFilterEnd(context map[string]interface{}) string {
	return Nprintf(`
data "google_billing_account" "account" {
  billing_account = "%{billing_acct}"
}

data "google_project" "project" {
}

resource "google_billing_budget" "budget" {
  billing_account = data.google_billing_account.account.id
  display_name = "Example Billing Budget%{random_suffix}"

  budget_filter {
    projects = []
  }

  amount {
    specified_amount {
      currency_code = "USD"
      units = "100000"
    }
  }

  threshold_rules {
    threshold_percent = 0.5
  }
  threshold_rules {
    threshold_percent = 0.9
    spend_basis = "FORECASTED_SPEND"
  }
}
`, context)
}

func TestBillingBudgetStateUpgradeV0(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		Attributes map[string]interface{}
		Expected   map[string]string
		Meta       interface{}
	}{
		"shorten long name": {
			Attributes: map[string]interface{}{
				"name": "billingAccounts/000000-111111-222222/budgets/9188612e-e4c0-4e69-9d14-9befebbcb87d",
			},
			Expected: map[string]string{
				"name": "9188612e-e4c0-4e69-9d14-9befebbcb87d",
			},
			Meta: &Config{},
		},
		"short name stays": {
			Attributes: map[string]interface{}{
				"name": "9188612e-e4c0-4e69-9d14-9befebbcb87d",
			},
			Expected: map[string]string{
				"name": "9188612e-e4c0-4e69-9d14-9befebbcb87d",
			},
			Meta: &Config{},
		},
	}
	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			actual, err := resourceBillingBudgetUpgradeV0(context.Background(), tc.Attributes, tc.Meta)

			if err != nil {
				t.Error(err)
			}

			for _, expectedName := range tc.Expected {
				if actual["name"] != expectedName {
					t.Errorf("expected: name -> %#v\n got: name -> %#v\n in: %#v",
						expectedName, actual["name"], actual)
				}
			}
		})
	}
}
