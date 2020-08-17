package google

import (
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
)

var IamProjectSchema = map[string]*schema.Schema{
	"project": {
		Type:             schema.TypeString,
		Optional:         true,
		Computed:         true,
		ForceNew:         true,
		DiffSuppressFunc: compareProjectName,
	},
}

// In google_project_iam_policy, project is required and not inferred by
// getProject.
var IamPolicyProjectSchema = map[string]*schema.Schema{
	"project": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		DiffSuppressFunc: compareProjectName,
	},
}

type ProjectIamUpdater struct {
	resourceId string
	Config     *Config
}

func NewProjectIamUpdater(d *schema.ResourceData, config *Config) (ResourceIamUpdater, error) {
	pid, err := getProject(d, config)
	if err != nil {
		return nil, err
	}

	d.Set("project", pid)

	return &ProjectIamUpdater{
		resourceId: pid,
		Config:     config,
	}, nil
}

// NewProjectIamPolicyUpdater is similar to NewProjectIamUpdater, except that it
// doesn't call getProject and only uses an explicitly set project.
func NewProjectIamPolicyUpdater(d *schema.ResourceData, config *Config) (ResourceIamUpdater, error) {
	return &ProjectIamUpdater{
		resourceId: d.Get("project").(string),
		Config:     config,
	}, nil
}

func ProjectIdParseFunc(d *schema.ResourceData, _ *Config) error {
	d.Set("project", d.Id())
	return nil
}

func (u *ProjectIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	projectId := GetResourceNameFromSelfLink(u.resourceId)
	p, err := u.Config.clientResourceManager.Projects.GetIamPolicy(projectId,
		&cloudresourcemanager.GetIamPolicyRequest{
			Options: &cloudresourcemanager.GetPolicyOptions{
				RequestedPolicyVersion: iamPolicyVersion,
			},
		}).Do()

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error retrieving IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return p, nil
}

func (u *ProjectIamUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	projectId := GetResourceNameFromSelfLink(u.resourceId)
	_, err := u.Config.clientResourceManager.Projects.SetIamPolicy(projectId,
		&cloudresourcemanager.SetIamPolicyRequest{
			Policy:     policy,
			UpdateMask: "bindings,etag,auditConfigs",
		}).Do()

	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error setting IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return nil
}

func (u *ProjectIamUpdater) GetResourceId() string {
	return u.resourceId
}

func (u *ProjectIamUpdater) GetMutexKey() string {
	return getProjectIamPolicyMutexKey(u.resourceId)
}

func (u *ProjectIamUpdater) DescribeResource() string {
	return fmt.Sprintf("project %q", u.resourceId)
}

func compareProjectName(_, old, new string, _ *schema.ResourceData) bool {
	// We can either get "projects/project-id" or "project-id", so strip any prefixes
	return GetResourceNameFromSelfLink(old) == GetResourceNameFromSelfLink(new)
}

func getProjectIamPolicyMutexKey(pid string) string {
	return fmt.Sprintf("iam-project-%s", pid)
}
