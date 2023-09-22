// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package certificatemanager

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleCertificateManagerCertificateMap() *schema.Resource {

	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceCertificateManagerCertificateMap().Schema)
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "name")
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceGoogleCertificateManagerCertificateMapRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleCertificateManagerCertificateMapRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	name := d.Get("name").(string)

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	id := fmt.Sprintf("projects/%s/locations/global/certificateMaps/%s", project, name)
	d.SetId(id)
	err = resourceCertificateManagerCertificateMapRead(d, meta)
	if err != nil {
		return err
	}

	if d.Id() == "" {
		return fmt.Errorf("%s not found", id)
	}
	return nil
}
