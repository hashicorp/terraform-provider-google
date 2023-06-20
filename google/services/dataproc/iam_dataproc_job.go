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

var IamDataprocJobSchema = map[string]*schema.Schema{
	"job_id": {
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

type DataprocJobIamUpdater struct {
	project string
	region  string
	jobId   string
	d       tpgresource.TerraformResourceData
	Config  *transport_tpg.Config
}

func NewDataprocJobUpdater(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (tpgiamresource.ResourceIamUpdater, error) {
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

	return &DataprocJobIamUpdater{
		project: project,
		region:  region,
		jobId:   d.Get("job_id").(string),
		d:       d,
		Config:  config,
	}, nil
}

func DataprocJobIdParseFunc(d *schema.ResourceData, config *transport_tpg.Config) error {
	fv, err := tpgresource.ParseRegionalFieldValue("jobs", d.Id(), "project", "region", "zone", d, config, true)
	if err != nil {
		return err
	}

	if err := d.Set("job_id", fv.Name); err != nil {
		return fmt.Errorf("Error setting job_id: %s", err)
	}
	if err := d.Set("project", fv.Project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("region", fv.Region); err != nil {
		return fmt.Errorf("Error setting region: %s", err)
	}

	// Explicitly set the id so imported resources have the same ID format as non-imported ones.
	d.SetId(fv.RelativeLink())
	return nil
}

func (u *DataprocJobIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	req := &dataproc.GetIamPolicyRequest{}

	userAgent, err := tpgresource.GenerateUserAgentString(u.d, u.Config.UserAgent)
	if err != nil {
		return nil, err
	}

	p, err := u.Config.NewDataprocClient(userAgent).Projects.Regions.Jobs.GetIamPolicy(u.GetResourceId(), req).Do()
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error retrieving IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	cloudResourcePolicy, err := dataprocToResourceManagerPolicy(p)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Invalid IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return cloudResourcePolicy, nil
}

func (u *DataprocJobIamUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	dataprocPolicy, err := resourceManagerToDataprocPolicy(policy)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Invalid IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	userAgent, err := tpgresource.GenerateUserAgentString(u.d, u.Config.UserAgent)
	if err != nil {
		return err
	}

	req := &dataproc.SetIamPolicyRequest{Policy: dataprocPolicy}
	_, err = u.Config.NewDataprocClient(userAgent).Projects.Regions.Jobs.SetIamPolicy(u.GetResourceId(), req).Do()
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error setting IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return nil
}

func (u *DataprocJobIamUpdater) GetResourceId() string {
	return fmt.Sprintf("projects/%s/regions/%s/jobs/%s", u.project, u.region, u.jobId)
}

func (u *DataprocJobIamUpdater) GetMutexKey() string {
	return fmt.Sprintf("iam-dataproc-job-%s-%s-%s", u.project, u.region, u.jobId)
}

func (u *DataprocJobIamUpdater) DescribeResource() string {
	return fmt.Sprintf("Dataproc Job %s/%s/%s", u.project, u.region, u.jobId)
}

func resourceManagerToDataprocPolicy(p *cloudresourcemanager.Policy) (*dataproc.Policy, error) {
	out := &dataproc.Policy{}
	err := tpgresource.Convert(p, out)
	if err != nil {
		return nil, errwrap.Wrapf("Cannot convert a dataproc policy to a cloudresourcemanager policy: {{err}}", err)
	}
	return out, nil
}

func dataprocToResourceManagerPolicy(p *dataproc.Policy) (*cloudresourcemanager.Policy, error) {
	out := &cloudresourcemanager.Policy{}
	err := tpgresource.Convert(p, out)
	if err != nil {
		return nil, errwrap.Wrapf("Cannot convert a cloudresourcemanager policy to a dataproc policy: {{err}}", err)
	}
	return out, nil
}
