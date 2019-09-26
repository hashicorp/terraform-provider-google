package google

import (
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
	Config  *Config
}

func NewDataprocClusterUpdater(d *schema.ResourceData, config *Config) (ResourceIamUpdater, error) {
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

	return &DataprocClusterIamUpdater{
		project: project,
		region:  region,
		cluster: d.Get("cluster").(string),
		Config:  config,
	}, nil
}

func DataprocClusterIdParseFunc(d *schema.ResourceData, config *Config) error {
	fv, err := parseRegionalFieldValue("clusters", d.Id(), "project", "region", "zone", d, config, true)
	if err != nil {
		return err
	}

	d.Set("project", fv.Project)
	d.Set("region", fv.Region)
	d.Set("cluster", fv.Name)

	// Explicitly set the id so imported resources have the same ID format as non-imported ones.
	d.SetId(fv.RelativeLink())
	return nil
}

func (u *DataprocClusterIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	req := &dataproc.GetIamPolicyRequest{}
	p, err := u.Config.clientDataproc.Projects.Regions.Clusters.GetIamPolicy(u.GetResourceId(), req).Do()
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

	req := &dataproc.SetIamPolicyRequest{Policy: dataprocPolicy}
	_, err = u.Config.clientDataproc.Projects.Regions.Clusters.SetIamPolicy(u.GetResourceId(), req).Do()
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
