// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package kms

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleKmsCryptoKey() *schema.Resource {
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceKMSCryptoKey().Schema)
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "name")
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "key_ring")

	return &schema.Resource{
		Read:   dataSourceGoogleKmsCryptoKeyRead,
		Schema: dsSchema,
	}

}

func dataSourceGoogleKmsCryptoKeyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	keyRingId, err := parseKmsKeyRingId(d.Get("key_ring").(string), config)
	if err != nil {
		return err
	}

	cryptoKeyId := KmsCryptoKeyId{
		KeyRingId: *keyRingId,
		Name:      d.Get("name").(string),
	}

	id := cryptoKeyId.CryptoKeyId()
	d.SetId(id)

	err = resourceKMSCryptoKeyRead(d, meta)
	if err != nil {
		return err
	}

	if err := tpgresource.SetDataSourceLabels(d); err != nil {
		return err
	}

	if d.Id() == "" {
		return fmt.Errorf("%s not found", id)
	}
	return nil
}
