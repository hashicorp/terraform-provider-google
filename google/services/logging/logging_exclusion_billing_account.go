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

var BillingAccountLoggingExclusionSchema = map[string]*schema.Schema{
	"billing_account": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
}

type BillingAccountLoggingExclusionUpdater struct {
	resourceType string
	resourceId   string
	userAgent    string
	Config       *transport_tpg.Config
}

func NewBillingAccountLoggingExclusionUpdater(d *schema.ResourceData, config *transport_tpg.Config) (ResourceLoggingExclusionUpdater, error) {
	billingAccount := d.Get("billing_account").(string)

	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return nil, err
	}

	return &BillingAccountLoggingExclusionUpdater{
		resourceType: "billingAccounts",
		resourceId:   billingAccount,
		userAgent:    userAgent,
		Config:       config,
	}, nil
}

func BillingAccountLoggingExclusionIdParseFunc(d *schema.ResourceData, _ *transport_tpg.Config) error {
	loggingExclusionId, err := ParseLoggingExclusionId(d.Id())
	if err != nil {
		return err
	}

	if "billingAccounts" != loggingExclusionId.resourceType {
		return fmt.Errorf("Error importing logging exclusion, invalid resourceType %#v", loggingExclusionId.resourceType)
	}

	if err := d.Set("billing_account", loggingExclusionId.ResourceId); err != nil {
		return fmt.Errorf("Error setting billing_account: %s", err)
	}
	return nil
}

func (u *BillingAccountLoggingExclusionUpdater) CreateLoggingExclusion(parent string, exclusion *logging.LogExclusion) error {
	_, err := u.Config.NewLoggingClient(u.userAgent).BillingAccounts.Exclusions.Create(parent, exclusion).Do()
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error creating logging exclusion for %s: {{err}}", u.DescribeResource()), err)
	}

	return nil
}

func (u *BillingAccountLoggingExclusionUpdater) ReadLoggingExclusion(id string) (*logging.LogExclusion, error) {
	exclusion, err := u.Config.NewLoggingClient(u.userAgent).BillingAccounts.Exclusions.Get(id).Do()

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error retrieving logging exclusion for %s: {{err}}", u.DescribeResource()), err)
	}

	return exclusion, nil
}

func (u *BillingAccountLoggingExclusionUpdater) UpdateLoggingExclusion(id string, exclusion *logging.LogExclusion, updateMask string) error {
	_, err := u.Config.NewLoggingClient(u.userAgent).BillingAccounts.Exclusions.Patch(id, exclusion).UpdateMask(updateMask).Do()
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error updating logging exclusion for %s: {{err}}", u.DescribeResource()), err)
	}

	return nil
}

func (u *BillingAccountLoggingExclusionUpdater) DeleteLoggingExclusion(id string) error {
	_, err := u.Config.NewLoggingClient(u.userAgent).BillingAccounts.Exclusions.Delete(id).Do()
	if err != nil {
		return errwrap.Wrap(fmt.Errorf("Error deleting logging exclusion for %s.", u.DescribeResource()), err)
	}

	return nil
}

func (u *BillingAccountLoggingExclusionUpdater) GetResourceType() string {
	return u.resourceType
}

func (u *BillingAccountLoggingExclusionUpdater) GetResourceId() string {
	return u.resourceId
}

func (u *BillingAccountLoggingExclusionUpdater) DescribeResource() string {
	return fmt.Sprintf("%q %q", u.resourceType, u.resourceId)
}
