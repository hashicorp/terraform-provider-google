package google

import (
	"fmt"
	"google.golang.org/api/bigtableadmin/v2"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
)

var IamBigtableInstanceSchema = map[string]*schema.Schema{
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
}

type BigtableInstanceIamUpdater struct {
	project  string
	instance string
	Config   *Config
}

func NewBigtableInstanceUpdater(d *schema.ResourceData, config *Config) (ResourceIamUpdater, error) {
	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}

	d.Set("project", project)

	return &BigtableInstanceIamUpdater{
		project:  project,
		instance: d.Get("instance").(string),
		Config:   config,
	}, nil
}

func BigtableInstanceIdParseFunc(d *schema.ResourceData, config *Config) error {
	fv, err := parseProjectFieldValue("instances", d.Id(), "project", d, config, false)
	if err != nil {
		return err
	}

	d.Set("project", fv.Project)
	d.Set("instance", fv.Name)

	// Explicitly set the id so imported resources have the same ID format as non-imported ones.
	d.SetId(fv.RelativeLink())
	return nil
}

func (u *BigtableInstanceIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	req := &bigtableadmin.GetIamPolicyRequest{}
	p, err := u.Config.clientBigtableProjectsInstances.GetIamPolicy(u.GetResourceId(), req).Do()
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error retrieving IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	cloudResourcePolicy, err := bigtableToResourceManagerPolicy(p)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Invalid IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return cloudResourcePolicy, nil
}

func (u *BigtableInstanceIamUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	bigtablePolicy, err := resourceManagerToBigtablePolicy(policy)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Invalid IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	req := &bigtableadmin.SetIamPolicyRequest{Policy: bigtablePolicy}
	_, err = u.Config.clientBigtableProjectsInstances.SetIamPolicy(u.GetResourceId(), req).Do()
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error setting IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return nil
}

func (u *BigtableInstanceIamUpdater) GetResourceId() string {
	return fmt.Sprintf("projects/%s/instances/%s", u.project, u.instance)
}

func (u *BigtableInstanceIamUpdater) GetMutexKey() string {
	return fmt.Sprintf("iam-bigtable-instance-%s-%s", u.project, u.instance)
}

func (u *BigtableInstanceIamUpdater) DescribeResource() string {
	return fmt.Sprintf("Bigtable Instance %s/%s", u.project, u.instance)
}

func resourceManagerToBigtablePolicy(p *cloudresourcemanager.Policy) (*bigtableadmin.Policy, error) {
	out := &bigtableadmin.Policy{}
	err := Convert(p, out)
	if err != nil {
		return nil, errwrap.Wrapf("Cannot convert a bigtable policy to a cloudresourcemanager policy: {{err}}", err)
	}
	return out, nil
}

func bigtableToResourceManagerPolicy(p *bigtableadmin.Policy) (*cloudresourcemanager.Policy, error) {
	out := &cloudresourcemanager.Policy{}
	err := Convert(p, out)
	if err != nil {
		return nil, errwrap.Wrapf("Cannot convert a cloudresourcemanager policy to a bigtable policy: {{err}}", err)
	}
	return out, nil
}
