// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package logging

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleLoggingProjectSettings() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleLoggingProjectSettingsRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The project for which to retrieve settings.`,
			},
			"disable_default_sink": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `If set to true, the _Default sink in newly created projects and folders will created in a disabled state. This can be used to automatically disable log storage if there is already an aggregated sink configured in the hierarchy. The _Default sink can be re-enabled manually if needed.`,
			},
			"kms_key_name": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The resource name for the configured Cloud KMS key.
				KMS key name format:
				"projects/[PROJECT_ID]/locations/[LOCATION]/keyRings/[KEYRING]/cryptoKeys/[KEY]"
				To enable CMEK for the bucket, set this field to a valid kmsKeyName for which the associated service account has the required cloudkms.cryptoKeyEncrypterDecrypter roles assigned for the key.
				The Cloud KMS key used by the bucket can be updated by changing the kmsKeyName to a new valid key name. Encryption operations that are in progress will be completed with the key that was in use when they started. Decryption operations will be completed using the key that was used at the time of encryption unless access to that key has been revoked.
				See [Enabling CMEK for Logging Buckets](https://cloud.google.com/logging/docs/routing/managed-encryption-storage) for more information.`,
			},
			"storage_location": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The storage location that Cloud Logging will use to create new resources when a location is needed but not explicitly provided.`,
			},
			"kms_service_account_id": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The service account associated with a project for which CMEK will apply.
				Before enabling CMEK for a logging bucket, you must first assign the cloudkms.cryptoKeyEncrypterDecrypter role to the service account associated with the project for which CMEK will apply. Use [v2.getCmekSettings](https://cloud.google.com/logging/docs/reference/v2/rest/v2/TopLevel/getCmekSettings#google.logging.v2.ConfigServiceV2.GetCmekSettings) to obtain the service account ID.
				See [Enabling CMEK for Logging Buckets](https://cloud.google.com/logging/docs/routing/managed-encryption-storage) for more information.`,
			},
			"logging_service_account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The service account for the given container. Sinks use this service account as their writerIdentity if no custom service account is provided.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource name of the CMEK settings.`,
			},
		},
	}
}

func dataSourceGoogleLoggingProjectSettingsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project := d.Get("project").(string)
	res, err := config.NewLoggingClient(userAgent).Projects.GetSettings(fmt.Sprintf("projects/%s", project)).Do()
	if err != nil {
		return transport_tpg.HandleDataSourceNotFoundError(err, d, fmt.Sprintf("LoggingProjectSettings %q", d.Id()), d.Id())
	}

	d.SetId(fmt.Sprintf("projects/%s/settings", project))

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading ProjectSettings: %s", err)
	}

	if err := d.Set("name", res.Name); err != nil {
		return fmt.Errorf("Error reading ProjectSettings: %s", err)
	}
	if err := d.Set("disable_default_sink", res.DisableDefaultSink); err != nil {
		return fmt.Errorf("Error reading ProjectSettings: %s", err)
	}
	if err := d.Set("kms_key_name", res.KmsKeyName); err != nil {
		return fmt.Errorf("Error reading ProjectSettings: %s", err)
	}
	if err := d.Set("storage_location", res.StorageLocation); err != nil {
		return fmt.Errorf("Error reading ProjectSettings: %s", err)
	}
	if err := d.Set("kms_service_account_id", res.KmsServiceAccountId); err != nil {
		return fmt.Errorf("Error reading ProjectSettings: %s", err)
	}
	if err := d.Set("logging_service_account_id", res.LoggingServiceAccountId); err != nil {
		return fmt.Errorf("Error reading ProjectSettings: %s", err)
	}
	return nil
}
