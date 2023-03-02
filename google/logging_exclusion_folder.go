package google

import (
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/logging/v2"
)

var FolderLoggingExclusionSchema = map[string]*schema.Schema{
	"folder": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		DiffSuppressFunc: optionalPrefixSuppress("folders/"),
	},
}

type FolderLoggingExclusionUpdater struct {
	resourceType string
	resourceId   string
	userAgent    string
	Config       *Config
}

func NewFolderLoggingExclusionUpdater(d *schema.ResourceData, config *Config) (ResourceLoggingExclusionUpdater, error) {
	folder := parseFolderId(d.Get("folder"))

	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return nil, err
	}

	return &FolderLoggingExclusionUpdater{
		resourceType: "folders",
		resourceId:   folder,
		userAgent:    userAgent,
		Config:       config,
	}, nil
}

func FolderLoggingExclusionIdParseFunc(d *schema.ResourceData, _ *Config) error {
	loggingExclusionId, err := parseLoggingExclusionId(d.Id())
	if err != nil {
		return err
	}

	if "folders" != loggingExclusionId.resourceType {
		return fmt.Errorf("Error importing logging exclusion, invalid resourceType %#v", loggingExclusionId.resourceType)
	}

	if err := d.Set("folder", loggingExclusionId.resourceId); err != nil {
		return fmt.Errorf("Error setting folder: %s", err)
	}
	return nil
}

func (u *FolderLoggingExclusionUpdater) CreateLoggingExclusion(parent string, exclusion *logging.LogExclusion) error {
	_, err := u.Config.NewLoggingClient(u.userAgent).Folders.Exclusions.Create(parent, exclusion).Do()
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error creating logging exclusion for %s: {{err}}", u.DescribeResource()), err)
	}

	return nil
}

func (u *FolderLoggingExclusionUpdater) ReadLoggingExclusion(id string) (*logging.LogExclusion, error) {
	exclusion, err := u.Config.NewLoggingClient(u.userAgent).Folders.Exclusions.Get(id).Do()

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error retrieving logging exclusion for %s: {{err}}", u.DescribeResource()), err)
	}

	return exclusion, nil
}

func (u *FolderLoggingExclusionUpdater) UpdateLoggingExclusion(id string, exclusion *logging.LogExclusion, updateMask string) error {
	_, err := u.Config.NewLoggingClient(u.userAgent).Folders.Exclusions.Patch(id, exclusion).UpdateMask(updateMask).Do()
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error updating logging exclusion for %s: {{err}}", u.DescribeResource()), err)
	}

	return nil
}

func (u *FolderLoggingExclusionUpdater) DeleteLoggingExclusion(id string) error {
	_, err := u.Config.NewLoggingClient(u.userAgent).Folders.Exclusions.Delete(id).Do()
	if err != nil {
		return errwrap.Wrap(fmt.Errorf("Error deleting logging exclusion for %s.", u.DescribeResource()), err)
	}

	return nil
}

func (u *FolderLoggingExclusionUpdater) GetResourceType() string {
	return u.resourceType
}

func (u *FolderLoggingExclusionUpdater) GetResourceId() string {
	return u.resourceId
}

func (u *FolderLoggingExclusionUpdater) DescribeResource() string {
	return fmt.Sprintf("%q %q", u.resourceType, u.resourceId)
}
