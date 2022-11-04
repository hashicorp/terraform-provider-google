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
    labels  = {
      label = "bar"
    }
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

func TestAccBillingBudget_billingBudgetUpdate(t *testing.T) {
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
				Config: testAccBillingBudget_billingBudgetUpdateStart(context),
			},
			{
				ResourceName:      "google_billing_budget.budget",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBillingBudget_billingBudgetUpdateRemoveFilter(context),
			},
			{
				ResourceName:      "google_billing_budget.budget",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBillingBudget_billingBudgetCalendarUpdate(context),
			},
			{
				ResourceName:      "google_billing_budget.budget",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBillingBudget_billingBudgetUpdate(context),
			},
			{
				ResourceName:      "google_billing_budget.budget",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBillingBudget_billingBudgetCustomPeriodUpdate(context),
			},
			{
				ResourceName:      "google_billing_budget.budget",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccBillingBudget_billingBudgetUpdateStart(context map[string]interface{}) string {
	return Nprintf(`
resource "google_pubsub_topic" "topic1" {
  name = "tf-test-billing-budget1-%{random_suffix}"
}
resource "google_pubsub_topic" "topic2" {
  name = "tf-test-billing-budget2-%{random_suffix}"
}

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
    labels  = {
      label = "bar"
    }
    credit_types_treatment = "EXCLUDE_ALL_CREDITS"
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

  all_updates_rule {
    pubsub_topic = google_pubsub_topic.topic1.id
  }
}
`, context)
}

func testAccBillingBudget_billingBudgetUpdateRemoveFilter(context map[string]interface{}) string {
	return Nprintf(`
resource "google_pubsub_topic" "topic1" {
  name = "tf-test-billing-budget1-%{random_suffix}"
}
resource "google_pubsub_topic" "topic2" {
  name = "tf-test-billing-budget2-%{random_suffix}"
}
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
    labels = {}
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

  all_updates_rule {
    pubsub_topic = google_pubsub_topic.topic1.id
  }
}
`, context)
}

func testAccBillingBudget_billingBudgetUpdate(context map[string]interface{}) string {
	return Nprintf(`
resource "google_pubsub_topic" "topic1" {
  name = "tf-test-billing-budget1-%{random_suffix}"
}
resource "google_pubsub_topic" "topic2" {
  name = "tf-test-billing-budget2-%{random_suffix}"
}
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
    labels  = {
      label1 = "bar2"
    }
    credit_types_treatment = "INCLUDE_ALL_CREDITS"
    services = ["services/24E6-581D-38E5"] # Bigquery
  }

  amount {
    specified_amount {
      currency_code = "USD"
      units = "2000"
    }
  }

  threshold_rules {
    threshold_percent = 0.5
  }
  threshold_rules {
    threshold_percent = 0.9
    spend_basis = "FORECASTED_SPEND"
  }

  all_updates_rule {
    pubsub_topic = google_pubsub_topic.topic2.id
  }
}
`, context)
}

func testAccBillingBudget_billingBudgetCalendarUpdate(context map[string]interface{}) string {
	return Nprintf(`
resource "google_pubsub_topic" "topic1" {
  name = "tf-test-billing-budget1-%{random_suffix}"
}
resource "google_pubsub_topic" "topic2" {
  name = "tf-test-billing-budget2-%{random_suffix}"
}
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
    labels  = {
      label1 = "bar2"
    }
	calendar_period = "YEAR"
  }

  amount {
    specified_amount {
      currency_code = "USD"
      units = "2000"
    }
  }

  threshold_rules {
    threshold_percent = 0.5
  }
  threshold_rules {
    threshold_percent = 0.9
  }

  all_updates_rule {
    pubsub_topic = google_pubsub_topic.topic2.id
  }
}
`, context)
}

func testAccBillingBudget_billingBudgetCustomPeriodUpdate(context map[string]interface{}) string {
	return Nprintf(`
resource "google_pubsub_topic" "topic1" {
  name = "tf-test-billing-budget1-%{random_suffix}"
}
resource "google_pubsub_topic" "topic2" {
  name = "tf-test-billing-budget2-%{random_suffix}"
}
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
    labels  = {
      label1 = "bar2"
    }
	custom_period {
	  start_date {
		year = 2022
		month = 1
		day = 1
	  }
	  end_date {
		year = 2023
		month = 12
		day = 31
	  }		
	}
  }

  amount {
    specified_amount {
      currency_code = "USD"
      units = "2000"
    }
  }

  threshold_rules {
    threshold_percent = 0.5
  }
  threshold_rules {
    threshold_percent = 0.9
  }

  all_updates_rule {
    pubsub_topic = google_pubsub_topic.topic2.id
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

func TestAccBillingBudget_budgetFilterProjectsOrdering(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org":             getTestOrgFromEnv(t),
		"billing_acct":    getTestBillingAccountFromEnv(t),
		"random_suffix_1": randString(t, 10),
		"random_suffix_2": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBillingBudgetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBillingBudget_budgetFilterProjectsOrdering1(context),
			},
			{
				ResourceName:      "google_billing_budget.budget",
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				Config:             testAccBillingBudget_budgetFilterProjectsOrdering2(context),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				ResourceName:      "google_billing_budget.budget",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccBillingBudget_budgetFilterProjectsOrdering1(context map[string]interface{}) string {
	return Nprintf(`

data "google_billing_account" "account" {
	billing_account = "%{billing_acct}"
}

resource "google_project" "project1" {
	project_id      = "tf-test-%{random_suffix_1}"
	name            = "tf-test-%{random_suffix_1}"
	org_id          = "%{org}"
	billing_account = data.google_billing_account.account.id
}

resource "google_project" "project2" {
	project_id      = "tf-test-%{random_suffix_2}"
	name            = "tf-test-%{random_suffix_2}"
	org_id          = "%{org}"
	billing_account = data.google_billing_account.account.id
}

resource "google_billing_budget" "budget" {
	billing_account = data.google_billing_account.account.id
	display_name    = "Example Billing Budget"

	budget_filter {
		projects = [
			"projects/${google_project.project1.number}",
			"projects/${google_project.project2.number}",
		]
	}

	amount {
		last_period_amount = true
	}

	threshold_rules {
		threshold_percent =  10.0
	}
}

`, context)
}

func testAccBillingBudget_budgetFilterProjectsOrdering2(context map[string]interface{}) string {
	return Nprintf(`

data "google_billing_account" "account" {
	billing_account = "%{billing_acct}"
}

resource "google_project" "project1" {
	project_id      = "tf-test-%{random_suffix_1}"
	name            = "tf-test-%{random_suffix_1}"
	org_id          = "%{org}"
	billing_account = data.google_billing_account.account.id
}

resource "google_project" "project2" {
	project_id      = "tf-test-%{random_suffix_2}"
	name            = "tf-test-%{random_suffix_2}"
	org_id          = "%{org}"
	billing_account = data.google_billing_account.account.id
}

resource "google_billing_budget" "budget" {
	billing_account = data.google_billing_account.account.id
	display_name    = "Example Billing Budget"

	budget_filter {
		projects = [
			"projects/${google_project.project2.number}",
			"projects/${google_project.project1.number}",
		]
	}

	amount {
		last_period_amount = true
	}

	threshold_rules {
		threshold_percent =  10.0
	}
}

`, context)
}
