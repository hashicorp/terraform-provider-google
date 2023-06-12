// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package bigtable

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-google/google/tpgiamresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/bigtableadmin/v2"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
)

var IamBigtableTableSchema = map[string]*schema.Schema{
	"instance": {
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
	"table": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
}

type BigtableTableIamUpdater struct {
	project  string
	instance string
	table    string
	d        tpgresource.TerraformResourceData
	Config   *transport_tpg.Config
}

func NewBigtableTableUpdater(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (tpgiamresource.ResourceIamUpdater, error) {
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return nil, err
	}

	if err := d.Set("project", project); err != nil {
		return nil, fmt.Errorf("Error setting project: %s", err)
	}

	return &BigtableTableIamUpdater{
		project:  project,
		instance: d.Get("instance").(string),
		table:    d.Get("table").(string),
		d:        d,
		Config:   config,
	}, nil
}

func BigtableTableIdParseFunc(d *schema.ResourceData, config *transport_tpg.Config) error {
	values := make(map[string]string)

	m, err := tpgresource.GetImportIdQualifiers([]string{"projects/(?P<project>[^/]+)/instances/(?P<instance>[^/]+)/tables/(?P<table>[^/]+)"}, d, config, d.Id())
	if err != nil {
		return err
	}

	project, _ := tpgresource.GetProject(d, config)

	for k, v := range m {
		values[k] = v
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}

	if err := d.Set("instance", values["instance"]); err != nil {
		return fmt.Errorf("Error setting instance: %s", err)
	}

	if err := d.Set("table", values["table"]); err != nil {
		return fmt.Errorf("Error setting table: %s", err)
	}

	// Explicitly set the id so imported resources have the same ID format as non-imported ones.
	d.SetId(fmt.Sprintf("projects/%s/instances/%s/tables/%s", project, values["instance"], values["table"]))
	return nil
}

func (u *BigtableTableIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	req := &bigtableadmin.GetIamPolicyRequest{}

	userAgent, err := tpgresource.GenerateUserAgentString(u.d, u.Config.UserAgent)
	if err != nil {
		return nil, err
	}

	p, err := u.Config.NewBigTableProjectsInstancesTablesClient(userAgent).GetIamPolicy(u.GetResourceId(), req).Do()
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error retrieving IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	cloudResourcePolicy, err := bigtableToResourceManagerPolicy(p)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Invalid IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return cloudResourcePolicy, nil
}

func (u *BigtableTableIamUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	bigtablePolicy, err := resourceManagerToBigtablePolicy(policy)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Invalid IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	req := &bigtableadmin.SetIamPolicyRequest{Policy: bigtablePolicy}

	userAgent, err := tpgresource.GenerateUserAgentString(u.d, u.Config.UserAgent)
	if err != nil {
		return err
	}

	_, err = u.Config.NewBigTableProjectsInstancesTablesClient(userAgent).SetIamPolicy(u.GetResourceId(), req).Do()
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error setting IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return nil
}

func (u *BigtableTableIamUpdater) GetResourceId() string {
	return fmt.Sprintf("projects/%s/instances/%s/tables/%s", u.project, u.instance, u.table)
}

func (u *BigtableTableIamUpdater) GetMutexKey() string {
	return fmt.Sprintf("iam-bigtable-instance-%s-%s-%s", u.project, u.instance, u.table)
}

func (u *BigtableTableIamUpdater) DescribeResource() string {
	return fmt.Sprintf("Bigtable Table %s/%s-%s", u.project, u.instance, u.table)
}
