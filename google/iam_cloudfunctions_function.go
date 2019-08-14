package google

import (
	"fmt"
	"google.golang.org/api/cloudfunctions/v1"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
)

var IamCloudFunctionsFunctionSchema = map[string]*schema.Schema{
	"project": {
		Type:     schema.TypeString,
		Computed: true,
		Optional: true,
		ForceNew: true,
	},
	"region": {
		Type:     schema.TypeString,
		Computed: true,
		Optional: true,
		ForceNew: true,
	},
	"function": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
}

type CloudFunctionsFunctionIamUpdater struct {
	project  string
	region   string
	function string
	Config *Config
}

func NewCloudFunctionsFunctionUpdater(d *schema.ResourceData, config *Config) (ResourceIamUpdater, error) {
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

	return &CloudFunctionsFunctionIamUpdater{
		project:  project,
		region:   region,
		function: d.Get("function").(string),
		Config:   config,
	}, nil
}

func CloudFunctionsFunctionIdParseFunc(d *schema.ResourceData, config *Config) error {
	fv, err := parseRegionalFieldValue("functions", d.Id(), "project", "region", "zone", d, config, true)
	if err != nil {
		return err
	}

	d.Set("project", fv.Project)
	d.Set("region", fv.Region)
	d.Set("function", fv.Name)

	funcId := &cloudFunctionId{
		Project: fv.Project,
		Region:  fv.Region,
		Name:    fv.Name,
	}
	// Explicitly set the id so imported resources have the same ID format as non-imported ones.
	d.SetId(funcId.cloudFunctionId())
	return nil
}

func (u *CloudFunctionsFunctionIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	p, err := u.Config.clientCloudFunctions.Projects.Locations.Functions.GetIamPolicy(u.GetResourceId()).Do()
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error retrieving IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	cloudResourcePolicy, err := cloudFunctionsToResourceManagerPolicy(p)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Invalid IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return cloudResourcePolicy, nil
}

func (u *CloudFunctionsFunctionIamUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	funcPolicy, err := resourceManagerToCloudFunctionsPolicy(policy)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Invalid IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	req := &cloudfunctions.SetIamPolicyRequest{Policy: funcPolicy}
	_, err = u.Config.clientCloudFunctions.Projects.Locations.Functions.SetIamPolicy(u.GetResourceId(), req).Do()
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error setting IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return nil
}

func (u *CloudFunctionsFunctionIamUpdater) GetResourceId() string {
	funcId := &cloudFunctionId{
		Project: u.project,
		Region:  u.region,
		Name:    u.function,
	}
	return funcId.cloudFunctionId()
}

func (u *CloudFunctionsFunctionIamUpdater) GetMutexKey() string {
	return fmt.Sprintf("iam-cloudfunctions-function-%s-%s-%s", u.project, u.region, u.function)
}

func (u *CloudFunctionsFunctionIamUpdater) DescribeResource() string {
	return fmt.Sprintf("Cloud Functions Function %s/%s/%s", u.project, u.region, u.function)
}

func resourceManagerToCloudFunctionsPolicy(p *cloudresourcemanager.Policy) (*cloudfunctions.Policy, error) {
	out := &cloudfunctions.Policy{}
	err := Convert(p, out)
	if err != nil {
		return nil, errwrap.Wrapf("Cannot convert a cloudfunctions policy to a cloudresourcemanager policy: {{err}}", err)
	}
	return out, nil
}

func cloudFunctionsToResourceManagerPolicy(p *cloudfunctions.Policy) (*cloudresourcemanager.Policy, error) {
	out := &cloudresourcemanager.Policy{}
	err := Convert(p, out)
	if err != nil {
		return nil, errwrap.Wrapf("Cannot convert a cloudresourcemanager policy to a cloudfunctions policy: {{err}}", err)
	}
	return out, nil
}
