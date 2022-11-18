package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGoogleLoggingProjectCmekSettings() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleLoggingProjectCmekSettingsRead,
		Schema: map[string]*schema.Schema{
			"kms_key_name": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `The resource name for the configured Cloud KMS key.
				KMS key name format:
				"projects/[PROJECT_ID]/locations/[LOCATION]/keyRings/[KEYRING]/cryptoKeys/[KEY]"
				To enable CMEK for the bucket, set this field to a valid kmsKeyName for which the associated service account has the required cloudkms.cryptoKeyEncrypterDecrypter roles assigned for the key.
				The Cloud KMS key used by the bucket can be updated by changing the kmsKeyName to a new valid key name. Encryption operations that are in progress will be completed with the key that was in use when they started. Decryption operations will be completed using the key that was used at the time of encryption unless access to that key has been revoked.
				See [Enabling CMEK for Logging Buckets](https://cloud.google.com/logging/docs/routing/managed-encryption-storage) for more information.`,
			},
			"kms_key_version_name": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The CryptoKeyVersion resource name for the configured Cloud KMS key.
				KMS key name format:
				"projects/[PROJECT_ID]/locations/[LOCATION]/keyRings/[KEYRING]/cryptoKeys/[KEY]/cryptoKeyVersions/[VERSION]"
				For example:
				"projects/my-project/locations/us-central1/keyRings/my-ring/cryptoKeys/my-key/cryptoKeyVersions/1"
				This is a read-only field used to convey the specific configured CryptoKeyVersion of kms_key that has been configured. It will be populated in cases where the CMEK settings are bound to a single key version.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource name of the CMEK settings.`,
			},
			"service_account_id": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The service account associated with a project for which CMEK will apply.
				Before enabling CMEK for a logging bucket, you must first assign the cloudkms.cryptoKeyEncrypterDecrypter role to the service account associated with the project for which CMEK will apply. Use [v2.getCmekSettings](https://cloud.google.com/logging/docs/reference/v2/rest/v2/TopLevel/getCmekSettings#google.logging.v2.ConfigServiceV2.GetCmekSettings) to obtain the service account ID.
				See [Enabling CMEK for Logging Buckets](https://cloud.google.com/logging/docs/routing/managed-encryption-storage) for more information.`,
			},
			"project": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceGoogleLoggingProjectCmekSettingsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "{{LoggingBasePath}}projects/{{project}}/cmekSettings")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for ProjectCmekSettings: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequest(config, "GET", billingProject, url, userAgent, nil)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("LoggingProjectCmekSettings %q", d.Id()))
	}

	d.SetId(fmt.Sprintf("projects/%s/cmekSettings", project))

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading ProjectCmekSettings: %s", err)
	}

	if err := d.Set("name", res["name"]); err != nil {
		return fmt.Errorf("Error reading ProjectCmekSettings: %s", err)
	}
	if err := d.Set("kms_key_name", res["kmsKeyName"]); err != nil {
		return fmt.Errorf("Error reading ProjectCmekSettings: %s", err)
	}
	if err := d.Set("kms_key_version_name", res["kmsKeyVersionName"]); err != nil {
		return fmt.Errorf("Error reading ProjectCmekSettings: %s", err)
	}
	if err := d.Set("service_account_id", res["serviceAccountId"]); err != nil {
		return fmt.Errorf("Error reading ProjectCmekSettings: %s", err)
	}

	return nil
}
