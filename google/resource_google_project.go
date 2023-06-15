// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"time"

	"github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/cloudresourcemanager/v1"
)

func PrefixedProject(pid string) string {
	return resourcemanager.PrefixedProject(pid)
}

func parseFolderId(v interface{}) string {
	return resourcemanager.ParseFolderId(v)
}

func EnableServiceUsageProjectServices(services []string, project, billingProject, userAgent string, config *transport_tpg.Config, timeout time.Duration) error {
	return resourcemanager.EnableServiceUsageProjectServices(services, project, billingProject, userAgent, config, timeout)
}

func ListCurrentlyEnabledServices(project, billingProject, userAgent string, config *transport_tpg.Config, timeout time.Duration) (map[string]struct{}, error) {
	return resourcemanager.ListCurrentlyEnabledServices(project, billingProject, userAgent, config, timeout)
}

// Retrieve the existing IAM Policy for a Project
func getProjectIamPolicy(project string, config *transport_tpg.Config) (*cloudresourcemanager.Policy, error) {
	return resourcemanager.GetProjectIamPolicy(project, config)
}
