// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
/*
	This file is a copy of mmv1/third_party/terraform/services/bigquery/iam_bigquery_dataset.go
	with new functions mergeAccess and GetCurrentResourceAccess
*/
package bigquery

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-google/google/tpgiamresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
)

var IamMemberBigqueryDatasetSchema = map[string]*schema.Schema{
	"dataset_id": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
	"project": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
		ForceNew: true,
	},
}

var bigqueryIamMemberAccessPrimitiveToRoleMap = map[string]string{
	"OWNER":  "roles/bigquery.dataOwner",
	"WRITER": "roles/bigquery.dataEditor",
	"READER": "roles/bigquery.dataViewer",
}

type BigqueryDatasetIamMemberUpdater struct {
	project   string
	datasetId string
	d         tpgresource.TerraformResourceData
	Config    *transport_tpg.Config
}

func NewBigqueryDatasetIamMemberUpdater(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (tpgiamresource.ResourceIamUpdater, error) {
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return nil, err
	}

	if err := d.Set("project", project); err != nil {
		return nil, fmt.Errorf("Error setting project: %s", err)
	}

	return &BigqueryDatasetIamMemberUpdater{
		project:   project,
		datasetId: d.Get("dataset_id").(string),
		d:         d,
		Config:    config,
	}, nil
}

func (u *BigqueryDatasetIamMemberUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	url := fmt.Sprintf("%s%s", u.Config.BigQueryBasePath, u.GetResourceId())

	userAgent, err := tpgresource.GenerateUserAgentString(u.d, u.Config.UserAgent)
	if err != nil {
		return nil, err
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    u.Config,
		Method:    "GET",
		Project:   u.project,
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error retrieving IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	policy, err := accessToPolicyForIamMember(res["access"])
	if err != nil {
		return nil, err
	}
	return policy, nil
}

func GetCurrentResourceAccess(u *BigqueryDatasetIamMemberUpdater) ([]interface{}, error) {
	url := fmt.Sprintf("%s%s", u.Config.BigQueryBasePath, u.GetResourceId())

	userAgent, err := tpgresource.GenerateUserAgentString(u.d, u.Config.UserAgent)
	if err != nil {
		return nil, err
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    u.Config,
		Method:    "GET",
		Project:   u.project,
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error retrieving IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	var access []interface{}
	if accessVal, ok := res["access"].([]interface{}); ok {
		access = accessVal
	} else if res["access"] == nil {
		access = []interface{}{} // Return an empty slice if the key is missing
	} else {
		return nil, fmt.Errorf("value under 'access' is not a slice of interface{}")
	}

	return access, nil
}

func mergeAccess(newAccess []map[string]interface{}, currAccess []interface{}) []map[string]interface{} {
	mergedAccess := make([]map[string]interface{}, 0, len(newAccess)+len(currAccess))
	mergedAccess = append(mergedAccess, newAccess...)

	for _, item := range currAccess {
		if itemMap, ok := item.(map[string]interface{}); ok {
			// Check if the item has a "dataset" key
			if _, ok := itemMap["dataset"]; ok {
				mergedAccess = append(mergedAccess, itemMap)
			}
		}
	}
	return mergedAccess
}

func (u *BigqueryDatasetIamMemberUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	url := fmt.Sprintf("%s%s", u.Config.BigQueryBasePath, u.GetResourceId())

	newAccess, err := policyToAccessForIamMember(policy)
	if err != nil {
		return err
	}

	currAccess, err := GetCurrentResourceAccess(u)
	if err != nil {
		return err
	}

	access := mergeAccess(newAccess, currAccess)

	obj := map[string]interface{}{
		"access": access,
	}

	userAgent, err := tpgresource.GenerateUserAgentString(u.d, u.Config.UserAgent)
	if err != nil {
		return err
	}

	_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    u.Config,
		Method:    "PATCH",
		Project:   u.project,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
	})
	if err != nil {
		return fmt.Errorf("Error creating DatasetAccess: %s", err)
	}

	return nil
}

func accessToPolicyForIamMember(access interface{}) (*cloudresourcemanager.Policy, error) {
	if access == nil {
		return nil, nil
	}
	roleToBinding := make(map[string]*cloudresourcemanager.Binding)

	accessArr := access.([]interface{})
	for _, v := range accessArr {
		memberRole := v.(map[string]interface{})
		rawRole, ok := memberRole["role"]
		if !ok {
			// "view" allows role to not be defined. It is a special dataset access construct, so ignore
			// If a user wants to manage "view" access they should use the `bigquery_dataset_access` resource
			continue
		}
		role := rawRole.(string)
		if iamRole, ok := bigqueryIamMemberAccessPrimitiveToRoleMap[role]; ok {
			// API changes certain IAM roles to legacy roles. Revert these changes
			role = iamRole
		}
		member, err := accessToIamMemberForIamMember(memberRole)
		if err != nil {
			return nil, err
		}
		// We have to combine bindings manually
		binding, ok := roleToBinding[role]
		if !ok {
			binding = &cloudresourcemanager.Binding{Role: role, Members: []string{}}
		}
		binding.Members = append(binding.Members, member)

		roleToBinding[role] = binding
	}
	bindings := make([]*cloudresourcemanager.Binding, 0)
	for _, v := range roleToBinding {
		bindings = append(bindings, v)
	}

	return &cloudresourcemanager.Policy{Bindings: bindings}, nil
}

func policyToAccessForIamMember(policy *cloudresourcemanager.Policy) ([]map[string]interface{}, error) {
	res := make([]map[string]interface{}, 0)
	if len(policy.AuditConfigs) != 0 {
		return nil, errors.New("Access policies not allowed on BigQuery Dataset IAM policies")
	}
	for _, binding := range policy.Bindings {
		if binding.Condition != nil {
			return nil, errors.New("IAM conditions not allowed on BigQuery Dataset IAM")
		}
		if fullRole, ok := bigqueryIamMemberAccessPrimitiveToRoleMap[binding.Role]; ok {
			return nil, fmt.Errorf("BigQuery Dataset legacy role %s is not allowed when using google_bigquery_dataset_iam resources. Please use the full form: %s", binding.Role, fullRole)
		}
		for _, member := range binding.Members {
			// Do not append any deleted members
			if strings.HasPrefix(member, "deleted:") {
				continue
			}
			access := map[string]interface{}{
				"role": binding.Role,
			}
			memberType, member, err := iamMemberToAccessForIamMember(member)
			if err != nil {
				return nil, err
			}
			access[memberType] = member
			res = append(res, access)
		}
	}

	return res, nil
}

// Returns the member access type and member for an IAM member.
// Dataset access uses different member types to identify groups, domains, etc.
// these types are used as keys in the access JSON payload
func iamMemberToAccessForIamMember(member string) (string, string, error) {
	if strings.HasPrefix(member, "deleted:") {
		return "", "", fmt.Errorf("BigQuery Dataset IAM member is deleted: %s", member)
	}
	pieces := strings.SplitN(member, ":", 2)
	if len(pieces) > 1 {
		switch pieces[0] {
		case "group":
			return "groupByEmail", pieces[1], nil
		case "domain":
			return "domain", pieces[1], nil
		case "iamMember":
			return "iamMember", pieces[1], nil
		case "user":
			return "userByEmail", pieces[1], nil
		case "serviceAccount":
			return "userByEmail", pieces[1], nil
		}
	}
	if member == "projectOwners" || member == "projectReaders" || member == "projectWriters" || member == "allAuthenticatedUsers" {
		// These are special BigQuery Dataset permissions
		return "specialGroup", member, nil
	}
	return "", "", fmt.Errorf("Failed to parse BigQuery Dataset IAM member type: %s", member)
}

func accessToIamMemberForIamMember(access map[string]interface{}) (string, error) {
	// One of the fields must be set, we have to find which IAM member type this newAccess to
	if member, ok := access["groupByEmail"]; ok {
		return fmt.Sprintf("group:%s", member.(string)), nil
	}
	if member, ok := access["domain"]; ok {
		return fmt.Sprintf("domain:%s", member.(string)), nil
	}
	if member, ok := access["specialGroup"]; ok {
		return member.(string), nil
	}
	if member, ok := access["iamMember"]; ok {
		return fmt.Sprintf("iamMember:%s", member.(string)), nil
	}
	if _, ok := access["view"]; ok {
		// view does not map to an IAM member, use access instead
		return "", fmt.Errorf("Failed to convert BigQuery Dataset access to IAM member. To use views with a dataset, please use dataset_access")
	}
	if _, ok := access["dataset"]; ok {
		// dataset does not map to an IAM member, use access instead
		return "", fmt.Errorf("Failed to convert BigQuery Dataset access to IAM member. To use views with a dataset, please use dataset_access")
	}
	if _, ok := access["routine"]; ok {
		// dataset does not map to an IAM member, use access instead
		return "", fmt.Errorf("Failed to convert BigQuery Dataset access to IAM member. To use views with a dataset, please use dataset_access")
	}
	if member, ok := access["userByEmail"]; ok {
		// service accounts have "gservice" in their email. This is best guess due to lost information
		if strings.Contains(member.(string), "gserviceaccount") {
			return fmt.Sprintf("serviceAccount:%s", member.(string)), nil
		}
		return fmt.Sprintf("user:%s", member.(string)), nil
	}
	return "", fmt.Errorf("Failed to identify IAM member from BigQuery Dataset access: %v", access)
}

func (u *BigqueryDatasetIamMemberUpdater) GetResourceId() string {
	return fmt.Sprintf("projects/%s/datasets/%s", u.project, u.datasetId)
}

// Matches the mutex of google_big_query_dataset_access
func (u *BigqueryDatasetIamMemberUpdater) GetMutexKey() string {
	return fmt.Sprintf("%s", u.datasetId)
}

func (u *BigqueryDatasetIamMemberUpdater) DescribeResource() string {
	return fmt.Sprintf("Bigquery Dataset %s/%s", u.project, u.datasetId)
}
