// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package monitoring_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/services/monitoring"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func TestAccMonitoringMonitoredProject_projectNumLongForm(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"project_id":    envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMonitoringMonitoredProjectDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringMonitoredProject_projectNumLongForm(context),
			},
			{
				ResourceName:            "google_monitoring_monitored_project.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"metrics_scope"},
			},
		},
	})
}

func TestAccMonitoringMonitoredProject_projectNumShortForm(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"project_id":    envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMonitoringMonitoredProjectDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringMonitoredProject_projectNumShortForm(context),
			},
			{
				ResourceName:            "google_monitoring_monitored_project.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"metrics_scope"},
			},
		},
	})
}

func testAccMonitoringMonitoredProject_projectNumLongForm(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_monitoring_monitored_project" "primary" {
  metrics_scope = "%{project_id}"
  name          = "locations/global/metricsScopes/%{project_id}/projects/${google_project.basic.number}"
}

resource "google_project" "basic" {
  project_id = "tf-test-m-id%{random_suffix}"
  name       = "tf-test-m-id%{random_suffix}-display"
  org_id     = "%{org_id}"
}
`, context)
}

func testAccMonitoringMonitoredProject_projectNumShortForm(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_monitoring_monitored_project" "primary" {
  metrics_scope = "%{project_id}"
  name          = "${google_project.basic.number}"
}

resource "google_project" "basic" {
  project_id = "tf-test-m-id%{random_suffix}"
  name       = "tf-test-m-id%{random_suffix}-display"
  org_id     = "%{org_id}"
}
`, context)
}

func TestUnitMonitoringMonitoredProject_nameDiffSuppress(t *testing.T) {
	for _, tc := range monitoringMonitoredProjectDiffSuppressTestCases {
		tc.Test(t)
	}
}

type MonitoringMonitoredProjectDiffSuppressTestCase struct {
	Name           string
	KeysToSuppress []string
	Before         map[string]interface{}
	After          map[string]interface{}
}

var monitoringMonitoredProjectDiffSuppressTestCases = []MonitoringMonitoredProjectDiffSuppressTestCase{
	// Project Id -> project Id
	{
		Name:           "short project id to long project id suppressed",
		KeysToSuppress: []string{"name"},
		Before: map[string]interface{}{
			"name": "sameId",
		},
		After: map[string]interface{}{
			"name": "locations/global/metricsScopes/projectId/projects/sameId",
		},
	},
	{
		Name:           "long project id to short project id suppressed",
		KeysToSuppress: []string{"name"},
		Before: map[string]interface{}{
			"name": "locations/global/metricsScopes/projectId/projects/sameId",
		},
		After: map[string]interface{}{
			"name": "sameId",
		},
	},
	{
		Name:           "short project id to long project id show diff",
		KeysToSuppress: []string{},
		Before: map[string]interface{}{
			"name": "oldId",
		},
		After: map[string]interface{}{
			"name": "locations/global/metricsScopes/projectId/projects/newId",
		},
	},
	{
		Name:           "long project id to short project id show diff",
		KeysToSuppress: []string{},
		Before: map[string]interface{}{
			"name": "locations/global/metricsScopes/projectId/projects/oldId",
		},
		After: map[string]interface{}{
			"name": "newId",
		},
	},

	// Project Num -> Project Num
	{
		Name:           "short project num to long project num suppressed",
		KeysToSuppress: []string{"name"},
		Before: map[string]interface{}{
			"name": "000000000000",
		},
		After: map[string]interface{}{
			"name": "locations/global/metricsScopes/projectId/projects/000000000000",
		},
	},
	{
		Name:           "long project num to short project num suppressed",
		KeysToSuppress: []string{"name"},
		Before: map[string]interface{}{
			"name": "locations/global/metricsScopes/projectId/projects/000000000000",
		},
		After: map[string]interface{}{
			"name": "000000000000",
		},
	},
	{
		Name:           "short project num to long project num show diff",
		KeysToSuppress: []string{},
		Before: map[string]interface{}{
			"name": "000000000000",
		},
		After: map[string]interface{}{
			"name": "locations/global/metricsScopes/projectId/projects/111111111111",
		},
	},
	{
		Name:           "long project num to short project num show diff",
		KeysToSuppress: []string{},
		Before: map[string]interface{}{
			"name": "locations/global/metricsScopes/projectId/projects/000000000000",
		},
		After: map[string]interface{}{
			"name": "111111111111",
		},
	},

	// Project id <--> Project num
	// Every variation of this should be suppressed. We cannot detect
	// if the project number matches the id within a diff suppress
	{
		Name:           "short project id to long project num suppressed",
		KeysToSuppress: []string{"name"},
		Before: map[string]interface{}{
			"name": "oldId",
		},
		After: map[string]interface{}{
			"name": "locations/global/metricsScopes/projectId/projects/111111111111",
		},
	},
	{
		Name:           "long project id to short project num suppressed",
		KeysToSuppress: []string{"name"},
		Before: map[string]interface{}{
			"name": "locations/global/metricsScopes/projectId/projects/oldId",
		},
		After: map[string]interface{}{
			"name": "111111111111",
		},
	},
	{
		Name:           "short project num to long project id suppressed",
		KeysToSuppress: []string{"name"},
		Before: map[string]interface{}{
			"name": "000000000000",
		},
		After: map[string]interface{}{
			"name": "locations/global/metricsScopes/projectId/projects/newId",
		},
	},
	{
		Name:           "long project num to short project id suppressed",
		KeysToSuppress: []string{"name"},
		Before: map[string]interface{}{
			"name": "locations/global/metricsScopes/projectId/projects/000000000000",
		},
		After: map[string]interface{}{
			"name": "newId",
		},
	},

	// Empty -> anything (resource creation)
	{
		Name:           "empty name to anything shows diff",
		KeysToSuppress: []string{},
		Before: map[string]interface{}{
			"name": "",
		},
		After: map[string]interface{}{
			"name": "newId",
		},
	},
}

func (tc *MonitoringMonitoredProjectDiffSuppressTestCase) Test(t *testing.T) {
	mockResourceDiff := &tpgresource.ResourceDiffMock{
		Before: tc.Before,
		After:  tc.After,
	}

	keySuppressionMap := map[string]bool{}
	for key := range tc.Before {
		keySuppressionMap[key] = false
	}
	for key := range tc.After {
		keySuppressionMap[key] = false
	}

	for _, key := range tc.KeysToSuppress {
		keySuppressionMap[key] = true
	}

	for key, tcSuppress := range keySuppressionMap {
		oldValue, ok := tc.Before[key]
		if !ok {
			oldValue = ""
		}
		newValue, ok := tc.After[key]
		if !ok {
			newValue = ""
		}
		suppressed := monitoring.ResourceMonitoringMonitoredProjectNameDiffSuppressFunc(key, fmt.Sprintf("%v", oldValue), fmt.Sprintf("%v", newValue), mockResourceDiff)
		if suppressed != tcSuppress {
			var expectation string
			if tcSuppress {
				expectation = "be"
			} else {
				expectation = "not be"
			}
			t.Errorf("Test %s: expected key `%s` to %s suppressed", tc.Name, key, expectation)
		}
	}
}
