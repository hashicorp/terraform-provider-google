package google

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
	spanner "google.golang.org/api/spanner/v1"
)

var IamSpannerInstanceSchema = map[string]*schema.Schema{
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

type SpannerInstanceIamUpdater struct {
	project  string
	instance string
	Config   *Config
}

func NewSpannerInstanceIamUpdater(d *schema.ResourceData, config *Config) (ResourceIamUpdater, error) {
	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}

	return &SpannerInstanceIamUpdater{
		project:  project,
		instance: d.Get("instance").(string),
		Config:   config,
	}, nil
}

func SpannerInstanceIdParseFunc(d *schema.ResourceData, config *Config) error {
	id, err := extractSpannerInstanceId(d.Id())
	if err != nil {
		return err
	}
	d.Set("instance", id.Instance)
	d.Set("project", id.Project)

	// Explicitly set the id so imported resources have the same ID format as non-imported ones.
	d.SetId(id.terraformId())
	return nil
}

func (u *SpannerInstanceIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	p, err := u.Config.clientSpanner.Projects.Instances.GetIamPolicy(spannerInstanceId{
		Project:  u.project,
		Instance: u.instance,
	}.instanceUri(), &spanner.GetIamPolicyRequest{}).Do()

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error retrieving IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	cloudResourcePolicy, err := spannerToResourceManagerPolicy(p)

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Invalid IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return cloudResourcePolicy, nil
}

func (u *SpannerInstanceIamUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	spannerPolicy, err := resourceManagerToSpannerPolicy(policy)

	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Invalid IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	_, err = u.Config.clientSpanner.Projects.Instances.SetIamPolicy(spannerInstanceId{
		Project:  u.project,
		Instance: u.instance,
	}.instanceUri(), &spanner.SetIamPolicyRequest{
		Policy: spannerPolicy,
	}).Do()

	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error setting IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return nil
}

func (u *SpannerInstanceIamUpdater) GetResourceId() string {
	return spannerInstanceId{
		Project:  u.project,
		Instance: u.instance,
	}.terraformId()
}

func (u *SpannerInstanceIamUpdater) GetMutexKey() string {
	return fmt.Sprintf("iam-spanner-instance-%s-%s", u.project, u.instance)
}

func (u *SpannerInstanceIamUpdater) DescribeResource() string {
	return fmt.Sprintf("Spanner Instance: %s/%s", u.project, u.instance)
}

type spannerInstanceId struct {
	Project  string
	Instance string
}

func (s spannerInstanceId) terraformId() string {
	return fmt.Sprintf("%s/%s", s.Project, s.Instance)
}

func (s spannerInstanceId) parentProjectUri() string {
	return fmt.Sprintf("projects/%s", s.Project)
}

func (s spannerInstanceId) instanceUri() string {
	return fmt.Sprintf("%s/instances/%s", s.parentProjectUri(), s.Instance)
}

func (s spannerInstanceId) instanceConfigUri(c string) string {
	return fmt.Sprintf("%s/instanceConfigs/%s", s.parentProjectUri(), c)
}

func extractSpannerInstanceId(id string) (*spannerInstanceId, error) {
	if !regexp.MustCompile("^" + ProjectRegex + "/[a-z0-9-]+$").Match([]byte(id)) {
		return nil, fmt.Errorf("Invalid spanner id format, expecting {projectId}/{instanceId}")
	}
	parts := strings.Split(id, "/")
	return &spannerInstanceId{
		Project:  parts[0],
		Instance: parts[1],
	}, nil
}
