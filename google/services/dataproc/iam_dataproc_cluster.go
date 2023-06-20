// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package dataproc

import (
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgiamresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/dataproc/v1"
)

var IamDataprocClusterSchema = map[string]*schema.Schema{
	"cluster": {
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
	"region": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
		ForceNew: true,
	},
}

type DataprocClusterIamUpdater struct {
	project string
	region  string
	cluster string
	d       tpgresource.TerraformResourceData
	Config  *transport_tpg.Config
}

func NewDataprocClusterUpdater(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (tpgiamresource.ResourceIamUpdater, error) {
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return nil, err
	}

	region, err := tpgresource.GetRegion(d, config)
	if err != nil {
		return nil, err
	}

	if err := d.Set("project", project); err != nil {
		return nil, fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("region", region); err != nil {
		return nil, fmt.Errorf("Error setting region: %s", err)
	}

	return &DataprocClusterIamUpdater{
		project: project,
		region:  region,
		cluster: d.Get("cluster").(string),
		d:       d,
		Config:  config,
	}, nil
}

func DataprocClusterIdParseFunc(d *schema.ResourceData, config *transport_tpg.Config) error {
	fv, err := tpgresource.ParseRegionalFieldValue("clusters", d.Id(), "project", "region", "zone", d, config, true)
	if err != nil {
		return err
	}

	if err := d.Set("project", fv.Project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("region", fv.Region); err != nil {
		return fmt.Errorf("Error setting region: %s", err)
	}
	if err := d.Set("cluster", fv.Name); err != nil {
		return fmt.Errorf("Error setting cluster: %s", err)
	}

	// Explicitly set the id so imported resources have the same ID format as non-imported ones.
	d.SetId(fv.RelativeLink())
	return nil
}

func (u *DataprocClusterIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	req := &dataproc.GetIamPolicyRequest{}

	userAgent, err := tpgresource.GenerateUserAgentString(u.d, u.Config.UserAgent)
	if err != nil {
		return nil, err
	}

	p, err := u.Config.NewDataprocClient(userAgent).Projects.Regions.Clusters.GetIamPolicy(u.GetResourceId(), req).Do()
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error retrieving IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	cloudResourcePolicy, err := dataprocToResourceManagerPolicy(p)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Invalid IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return cloudResourcePolicy, nil
}

func (u *DataprocClusterIamUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	dataprocPolicy, err := resourceManagerToDataprocPolicy(policy)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Invalid IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	userAgent, err := tpgresource.GenerateUserAgentString(u.d, u.Config.UserAgent)
	if err != nil {
		return err
	}

	req := &dataproc.SetIamPolicyRequest{Policy: dataprocPolicy}
	_, err = u.Config.NewDataprocClient(userAgent).Projects.Regions.Clusters.SetIamPolicy(u.GetResourceId(), req).Do()
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error setting IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return nil
}

func (u *DataprocClusterIamUpdater) GetResourceId() string {
	return fmt.Sprintf("projects/%s/regions/%s/clusters/%s", u.project, u.region, u.cluster)
}

func (u *DataprocClusterIamUpdater) GetMutexKey() string {
	return fmt.Sprintf("iam-dataproc-cluster-%s-%s-%s", u.project, u.region, u.cluster)
}

func (u *DataprocClusterIamUpdater) DescribeResource() string {
	return fmt.Sprintf("Dataproc Cluster %s/%s/%s", u.project, u.region, u.cluster)
}
