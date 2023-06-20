// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package logging

import (
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/logging/v2"
)

var ProjectLoggingExclusionSchema = map[string]*schema.Schema{
	"project": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
		ForceNew: true,
	},
}

type ProjectLoggingExclusionUpdater struct {
	resourceType string
	resourceId   string
	userAgent    string
	Config       *transport_tpg.Config
}

func NewProjectLoggingExclusionUpdater(d *schema.ResourceData, config *transport_tpg.Config) (ResourceLoggingExclusionUpdater, error) {
	pid, err := tpgresource.GetProject(d, config)
	if err != nil {
		return nil, err
	}

	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return nil, err
	}

	return &ProjectLoggingExclusionUpdater{
		resourceType: "projects",
		resourceId:   pid,
		userAgent:    userAgent,
		Config:       config,
	}, nil
}

func ProjectLoggingExclusionIdParseFunc(d *schema.ResourceData, config *transport_tpg.Config) error {
	loggingExclusionId, err := ParseLoggingExclusionId(d.Id())
	if err != nil {
		return err
	}

	if "projects" != loggingExclusionId.resourceType {
		return fmt.Errorf("Error importing logging exclusion, invalid resourceType %#v", loggingExclusionId.resourceType)
	}

	if config.Project != loggingExclusionId.ResourceId {
		if err := d.Set("project", loggingExclusionId.ResourceId); err != nil {
			return fmt.Errorf("Error setting project: %s", err)
		}
	}

	return nil
}

func (u *ProjectLoggingExclusionUpdater) CreateLoggingExclusion(parent string, exclusion *logging.LogExclusion) error {
	_, err := u.Config.NewLoggingClient(u.userAgent).Projects.Exclusions.Create(parent, exclusion).Do()
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error creating logging exclusion for %s: {{err}}", u.DescribeResource()), err)
	}

	return nil
}

func (u *ProjectLoggingExclusionUpdater) ReadLoggingExclusion(id string) (*logging.LogExclusion, error) {
	exclusion, err := u.Config.NewLoggingClient(u.userAgent).Projects.Exclusions.Get(id).Do()

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error retrieving logging exclusion for %s: {{err}}", u.DescribeResource()), err)
	}

	return exclusion, nil
}

func (u *ProjectLoggingExclusionUpdater) UpdateLoggingExclusion(id string, exclusion *logging.LogExclusion, updateMask string) error {
	_, err := u.Config.NewLoggingClient(u.userAgent).Projects.Exclusions.Patch(id, exclusion).UpdateMask(updateMask).Do()
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error updating logging exclusion for %s: {{err}}", u.DescribeResource()), err)
	}

	return nil
}

func (u *ProjectLoggingExclusionUpdater) DeleteLoggingExclusion(id string) error {
	_, err := u.Config.NewLoggingClient(u.userAgent).Projects.Exclusions.Delete(id).Do()
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error deleting logging exclusion for %s: {{err}}", u.DescribeResource()), err)
	}

	return nil
}

func (u *ProjectLoggingExclusionUpdater) GetResourceType() string {
	return u.resourceType
}

func (u *ProjectLoggingExclusionUpdater) GetResourceId() string {
	return u.resourceId
}

func (u *ProjectLoggingExclusionUpdater) DescribeResource() string {
	return fmt.Sprintf("%q %q", u.resourceType, u.resourceId)
}
