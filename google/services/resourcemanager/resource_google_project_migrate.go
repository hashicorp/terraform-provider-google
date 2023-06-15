// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/tpgiamresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/cloudresourcemanager/v1"
)

func resourceGoogleProjectMigrateState(v int, s *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	if s.Empty() {
		log.Println("[DEBUG] Empty InstanceState; nothing to migrate.")
		return s, nil
	}

	switch v {
	case 0:
		log.Println("[INFO] Found Google Project State v0; migrating to v1")
		s, err := migrateGoogleProjectStateV0toV1(s, meta.(*transport_tpg.Config))
		if err != nil {
			return s, err
		}
		return s, nil
	default:
		return s, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

// This migration adjusts google_project resources to include several additional attributes
// required to support project creation/deletion that was added in V1.
func migrateGoogleProjectStateV0toV1(s *terraform.InstanceState, config *transport_tpg.Config) (*terraform.InstanceState, error) {
	log.Printf("[DEBUG] Attributes before migration: %#v", s.Attributes)

	s.Attributes["skip_delete"] = "true"
	s.Attributes["project_id"] = s.ID

	if s.Attributes["policy_data"] != "" {
		p, err := GetProjectIamPolicy(s.ID, config)
		if err != nil {
			return s, fmt.Errorf("Could not retrieve project's IAM policy while attempting to migrate state from V0 to V1: %v", err)
		}
		s.Attributes["policy_etag"] = p.Etag
	}

	log.Printf("[DEBUG] Attributes after migration: %#v", s.Attributes)
	return s, nil
}

// Retrieve the existing IAM Policy for a Project
func GetProjectIamPolicy(project string, config *transport_tpg.Config) (*cloudresourcemanager.Policy, error) {
	p, err := config.NewResourceManagerClient(config.UserAgent).Projects.GetIamPolicy(project,
		&cloudresourcemanager.GetIamPolicyRequest{
			Options: &cloudresourcemanager.GetPolicyOptions{
				RequestedPolicyVersion: tpgiamresource.IamPolicyVersion,
			},
		}).Do()

	if err != nil {
		return nil, fmt.Errorf("Error retrieving IAM policy for project %q: %s", project, err)
	}
	return p, nil
}
