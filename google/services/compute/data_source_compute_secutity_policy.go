// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleComputeSecurityPolicy() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceComputeSecurityPolicy().Schema)

	// Set 'Optional' schema elements
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "name")
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project")
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "self_link")

	return &schema.Resource{
		Read:   dataSourceComputSecurityPolicyRead,
		Schema: dsSchema,
	}
}

func dataSourceComputSecurityPolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	id := ""

	if name, ok := d.GetOk("name"); ok {
		project, err := tpgresource.GetProject(d, config)
		if err != nil {
			return err
		}

		id = fmt.Sprintf("projects/%s/global/securityPolicies/%s", project, name.(string))
		d.SetId(id)
	} else if selfLink, ok := d.GetOk("self_link"); ok {
		parsed, err := tpgresource.ParseSecurityPolicyFieldValue(selfLink.(string), d, config)
		if err != nil {
			return err
		}

		if err := d.Set("name", parsed.Name); err != nil {
			return fmt.Errorf("Error setting name: %s", err)
		}

		if err := d.Set("project", parsed.Project); err != nil {
			return fmt.Errorf("Error setting project: %s", err)
		}

		id = fmt.Sprintf("projects/%s/global/securityPolicies/%s", parsed.Project, parsed.Name)
		d.SetId(id)
	} else {
		return errors.New("Must provide either `self_link` or `name`")
	}

	err := resourceComputeSecurityPolicyRead(d, meta)
	if err != nil {
		return err
	}

	if d.Id() == "" {
		return fmt.Errorf("%s not found", id)
	}

	return nil
}
