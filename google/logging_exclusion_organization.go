package google

import (
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/logging/v2"
)

var OrganizationLoggingExclusionSchema = map[string]*schema.Schema{
	"org_id": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
}

type OrganizationLoggingExclusionUpdater struct {
	resourceType string
	resourceId   string
	userAgent    string
	Config       *Config
}

func NewOrganizationLoggingExclusionUpdater(d *schema.ResourceData, config *Config) (ResourceLoggingExclusionUpdater, error) {
	organization := d.Get("org_id").(string)

	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return nil, err
	}

	return &OrganizationLoggingExclusionUpdater{
		resourceType: "organizations",
		resourceId:   organization,
		userAgent:    userAgent,
		Config:       config,
	}, nil
}

func organizationLoggingExclusionIdParseFunc(d *schema.ResourceData, _ *Config) error {
	loggingExclusionId, err := parseLoggingExclusionId(d.Id())
	if err != nil {
		return err
	}

	if "organizations" != loggingExclusionId.resourceType {
		return fmt.Errorf("Error importing logging exclusion, invalid resourceType %#v", loggingExclusionId.resourceType)
	}

	if err := d.Set("org_id", loggingExclusionId.resourceId); err != nil {
		return fmt.Errorf("Error setting org_id: %s", err)
	}
	return nil
}

func (u *OrganizationLoggingExclusionUpdater) CreateLoggingExclusion(parent string, exclusion *logging.LogExclusion) error {
	_, err := u.Config.NewLoggingClient(u.userAgent).Organizations.Exclusions.Create(parent, exclusion).Do()
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error creating logging exclusion for %s: {{err}}", u.DescribeResource()), err)
	}

	return nil
}

func (u *OrganizationLoggingExclusionUpdater) ReadLoggingExclusion(id string) (*logging.LogExclusion, error) {
	exclusion, err := u.Config.NewLoggingClient(u.userAgent).Organizations.Exclusions.Get(id).Do()

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error retrieving logging exclusion for %s: {{err}}", u.DescribeResource()), err)
	}

	return exclusion, nil
}

func (u *OrganizationLoggingExclusionUpdater) UpdateLoggingExclusion(id string, exclusion *logging.LogExclusion, updateMask string) error {
	_, err := u.Config.NewLoggingClient(u.userAgent).Organizations.Exclusions.Patch(id, exclusion).UpdateMask(updateMask).Do()
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error updating logging exclusion for %s: {{err}}", u.DescribeResource()), err)
	}

	return nil
}

func (u *OrganizationLoggingExclusionUpdater) DeleteLoggingExclusion(id string) error {
	_, err := u.Config.NewLoggingClient(u.userAgent).Organizations.Exclusions.Delete(id).Do()
	if err != nil {
		return errwrap.Wrap(fmt.Errorf("Error deleting logging exclusion for %s.", u.DescribeResource()), err)
	}

	return nil
}

func (u *OrganizationLoggingExclusionUpdater) GetResourceType() string {
	return u.resourceType
}

func (u *OrganizationLoggingExclusionUpdater) GetResourceId() string {
	return u.resourceId
}

func (u *OrganizationLoggingExclusionUpdater) DescribeResource() string {
	return fmt.Sprintf("%q %q", u.resourceType, u.resourceId)
}
