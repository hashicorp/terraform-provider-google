// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"google.golang.org/api/iam/v1"
)

func ResourceGoogleServiceAccountKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceGoogleServiceAccountKeyCreate,
		Read:   resourceGoogleServiceAccountKeyRead,
		Delete: resourceGoogleServiceAccountKeyDelete,
		Schema: map[string]*schema.Schema{
			// Required
			"service_account_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the parent service account of the key. This can be a string in the format {ACCOUNT} or projects/{PROJECT_ID}/serviceAccounts/{ACCOUNT}, where {ACCOUNT} is the email address or unique id of the service account. If the {ACCOUNT} syntax is used, the project will be inferred from the provider's configuration.`,
			},
			// Optional
			"key_algorithm": {
				Type:         schema.TypeString,
				Default:      "KEY_ALG_RSA_2048",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"KEY_ALG_UNSPECIFIED", "KEY_ALG_RSA_1024", "KEY_ALG_RSA_2048"}, false),
				Description:  `The algorithm used to generate the key, used only on create. KEY_ALG_RSA_2048 is the default algorithm. Valid values are: "KEY_ALG_RSA_1024", "KEY_ALG_RSA_2048".`,
			},
			"private_key_type": {
				Type:         schema.TypeString,
				Default:      "TYPE_GOOGLE_CREDENTIALS_FILE",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"TYPE_UNSPECIFIED", "TYPE_PKCS12_FILE", "TYPE_GOOGLE_CREDENTIALS_FILE"}, false),
			},
			"public_key_type": {
				Type:         schema.TypeString,
				Default:      "TYPE_X509_PEM_FILE",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"TYPE_NONE", "TYPE_X509_PEM_FILE", "TYPE_RAW_PUBLIC_KEY"}, false),
			},
			"public_key_data": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"key_algorithm", "private_key_type"},
				Description:   `A field that allows clients to upload their own public key. If set, use this public key data to create a service account key for given service account. Please note, the expected format for this field is a base64 encoded X509_PEM.`,
			},
			"keepers": {
				Description: "Arbitrary map of values that, when changed, will trigger recreation of resource.",
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
			},
			// Computed
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Description: `The name used for this key pair`,
			},
			"public_key": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Description: `The public key, base64 encoded`,
			},
			"private_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: `The private key in JSON format, base64 encoded. This is what you normally get as a file when creating service account keys through the CLI or web console. This is only populated when creating a new key.`,
			},
			"valid_after": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The key can be used after this timestamp. A timestamp in RFC3339 UTC "Zulu" format, accurate to nanoseconds. Example: "2014-10-02T15:01:23.045123456Z".`,
			},
			"valid_before": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The key can be used before this timestamp. A timestamp in RFC3339 UTC "Zulu" format, accurate to nanoseconds. Example: "2014-10-02T15:01:23.045123456Z".`,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceGoogleServiceAccountKeyCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	serviceAccountName, err := tpgresource.ServiceAccountFQN(d.Get("service_account_id").(string), d, config)
	if err != nil {
		return err
	}

	var sak *iam.ServiceAccountKey

	if d.Get("public_key_data").(string) != "" {
		ru := &iam.UploadServiceAccountKeyRequest{
			PublicKeyData: d.Get("public_key_data").(string),
		}
		sak, err = config.NewIamClient(userAgent).Projects.ServiceAccounts.Keys.Upload(serviceAccountName, ru).Do()
		if err != nil {
			return fmt.Errorf("Error creating service account key: %s", err)
		}
	} else {
		rc := &iam.CreateServiceAccountKeyRequest{
			KeyAlgorithm:   d.Get("key_algorithm").(string),
			PrivateKeyType: d.Get("private_key_type").(string),
		}
		sak, err = config.NewIamClient(userAgent).Projects.ServiceAccounts.Keys.Create(serviceAccountName, rc).Do()
		if err != nil {
			return fmt.Errorf("Error creating service account key: %s", err)
		}
	}

	d.SetId(sak.Name)
	// Data only available on create.
	if err := d.Set("valid_after", sak.ValidAfterTime); err != nil {
		return fmt.Errorf("Error setting valid_after: %s", err)
	}
	if err := d.Set("valid_before", sak.ValidBeforeTime); err != nil {
		return fmt.Errorf("Error setting valid_before: %s", err)
	}
	if err := d.Set("private_key", sak.PrivateKeyData); err != nil {
		return fmt.Errorf("Error setting private_key: %s", err)
	}

	err = ServiceAccountKeyWaitTime(config.NewIamClient(userAgent).Projects.ServiceAccounts.Keys, d.Id(), d.Get("public_key_type").(string), "Creating Service account key", 4*time.Minute)
	if err != nil {
		return err
	}
	return resourceGoogleServiceAccountKeyRead(d, meta)
}

func resourceGoogleServiceAccountKeyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	publicKeyType := d.Get("public_key_type").(string)

	// Confirm the service account key exists
	sak, err := config.NewIamClient(userAgent).Projects.ServiceAccounts.Keys.Get(d.Id()).PublicKeyType(publicKeyType).Do()
	if err != nil {
		if err = transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Service Account Key %q", d.Id())); err == nil {
			return nil
		} else {
			// This resource also returns 403 when it's not found.
			if transport_tpg.IsGoogleApiErrorWithCode(err, 403) {
				log.Printf("[DEBUG] Got a 403 error trying to read service account key %s, assuming it's gone.", d.Id())
				d.SetId("")
				return nil
			} else {
				return err
			}
		}
	}

	if err := d.Set("name", sak.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("key_algorithm", sak.KeyAlgorithm); err != nil {
		return fmt.Errorf("Error setting key_algorithm: %s", err)
	}
	if err := d.Set("public_key", sak.PublicKeyData); err != nil {
		return fmt.Errorf("Error setting public_key: %s", err)
	}
	return nil
}

func resourceGoogleServiceAccountKeyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	_, err = config.NewIamClient(userAgent).Projects.ServiceAccounts.Keys.Delete(d.Id()).Do()

	if err != nil {
		if err = transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Service Account Key %q", d.Id())); err == nil {
			return nil
		} else {
			// This resource also returns 403 when it's not found.
			if transport_tpg.IsGoogleApiErrorWithCode(err, 403) {
				log.Printf("[DEBUG] Got a 403 error trying to read service account key %s, assuming it's gone.", d.Id())
				d.SetId("")
				return nil
			} else {
				return err
			}
		}
	}

	d.SetId("")
	return nil
}
