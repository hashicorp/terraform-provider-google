package google

import (
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
	Config  *Config
}

func NewDataprocJobUpdater(d *schema.ResourceData, config *Config) (ResourceIamUpdater, error) {
	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return nil, err
	}

	d.Set("project", project)
	d.Set("region", region)

	return &DataprocJobIamUpdater{
		project: project,
		region:  region,
		jobId:   d.Get("job_id").(string),
		Config:  config,
	}, nil
}

func DataprocJobIdParseFunc(d *schema.ResourceData, config *Config) error {
	fv, err := parseRegionalFieldValue("jobs", d.Id(), "project", "region", "zone", d, config, true)
	if err != nil {
		return err
	}

	d.Set("job_id", fv.Name)
	d.Set("project", fv.Project)
	d.Set("region", fv.Region)

	// Explicitly set the id so imported resources have the same ID format as non-imported ones.
	d.SetId(fv.RelativeLink())
	return nil
}

func (u *DataprocJobIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	req := &dataproc.GetIamPolicyRequest{}
	p, err := u.Config.clientDataproc.Projects.Regions.Jobs.GetIamPolicy(u.GetResourceId(), req).Do()
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

	req := &dataproc.SetIamPolicyRequest{Policy: dataprocPolicy}
	_, err = u.Config.clientDataproc.Projects.Regions.Jobs.SetIamPolicy(u.GetResourceId(), req).Do()
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
	err := Convert(p, out)
	if err != nil {
		return nil, errwrap.Wrapf("Cannot convert a dataproc policy to a cloudresourcemanager policy: {{err}}", err)
	}
	return out, nil
}

func dataprocToResourceManagerPolicy(p *dataproc.Policy) (*cloudresourcemanager.Policy, error) {
	out := &cloudresourcemanager.Policy{}
	err := Convert(p, out)
	if err != nil {
		return nil, errwrap.Wrapf("Cannot convert a cloudresourcemanager policy to a dataproc policy: {{err}}", err)
	}
	return out, nil
}
