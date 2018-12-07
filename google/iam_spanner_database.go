package google

import (
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/spanner/v1"
)

var IamSpannerDatabaseSchema = map[string]*schema.Schema{
	"instance": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
	"database": {
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

type SpannerDatabaseIamUpdater struct {
	project  string
	instance string
	database string
	Config   *Config
}

func NewSpannerDatabaseIamUpdater(d *schema.ResourceData, config *Config) (ResourceIamUpdater, error) {
	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}

	return &SpannerDatabaseIamUpdater{
		project:  project,
		instance: d.Get("instance").(string),
		database: d.Get("database").(string),
		Config:   config,
	}, nil
}

func SpannerDatabaseIdParseFunc(d *schema.ResourceData, config *Config) error {
	_, err := resourceSpannerDatabaseImport("database")(d, config)
	return err
}

func (u *SpannerDatabaseIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	p, err := u.Config.clientSpanner.Projects.Instances.Databases.GetIamPolicy(spannerDatabaseId{
		Project:  u.project,
		Database: u.database,
		Instance: u.instance,
	}.databaseUri(), &spanner.GetIamPolicyRequest{}).Do()

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error retrieving IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	cloudResourcePolicy, err := spannerToResourceManagerPolicy(p)

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Invalid IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return cloudResourcePolicy, nil
}

func (u *SpannerDatabaseIamUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	spannerPolicy, err := resourceManagerToSpannerPolicy(policy)

	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Invalid IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	_, err = u.Config.clientSpanner.Projects.Instances.Databases.SetIamPolicy(spannerDatabaseId{
		Project:  u.project,
		Database: u.database,
		Instance: u.instance,
	}.databaseUri(), &spanner.SetIamPolicyRequest{
		Policy: spannerPolicy,
	}).Do()

	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error setting IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return nil
}

func (u *SpannerDatabaseIamUpdater) GetResourceId() string {
	return spannerDatabaseId{
		Project:  u.project,
		Instance: u.instance,
		Database: u.database,
	}.terraformId()
}

func (u *SpannerDatabaseIamUpdater) GetMutexKey() string {
	return fmt.Sprintf("iam-spanner-database-%s-%s-%s", u.project, u.instance, u.database)
}

func (u *SpannerDatabaseIamUpdater) DescribeResource() string {
	return fmt.Sprintf("Spanner Database: %s/%s/%s", u.project, u.instance, u.database)
}

func resourceManagerToSpannerPolicy(p *cloudresourcemanager.Policy) (*spanner.Policy, error) {
	out := &spanner.Policy{}
	err := Convert(p, out)
	if err != nil {
		return nil, errwrap.Wrapf("Cannot convert a resourcemanager policy to a spanner policy: {{err}}", err)
	}
	return out, nil
}

func spannerToResourceManagerPolicy(p *spanner.Policy) (*cloudresourcemanager.Policy, error) {
	out := &cloudresourcemanager.Policy{}
	err := Convert(p, out)
	if err != nil {
		return nil, errwrap.Wrapf("Cannot convert a spanner policy to a resourcemanager policy: {{err}}", err)
	}
	return out, nil
}
