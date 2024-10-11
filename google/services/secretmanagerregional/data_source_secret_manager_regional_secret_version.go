// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package secretmanagerregional

import (
	"encoding/base64"
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceSecretManagerRegionalRegionalSecretVersion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSecretManagerRegionalRegionalSecretVersionRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"location": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"secret": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"destroy_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"secret_data": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"customer_managed_encryption": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kms_key_version_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"is_secret_data_base64": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceSecretManagerRegionalRegionalSecretVersionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	secretRegex := regexp.MustCompile("projects/(.+)/locations/(.+)/secrets/(.+)$")
	parts := secretRegex.FindStringSubmatch(d.Get("secret").(string))

	var project string

	// if reference of the secret is provided in the secret field
	if len(parts) == 4 {
		// Store values of project to set in state
		project = parts[1]
		if d.Get("project").(string) != "" && d.Get("project").(string) != parts[1] {
			return fmt.Errorf("The project set on this secret version (%s) is not equal to the project where this secret exists (%s).", d.Get("project").(string), parts[1])
		}
		if d.Get("location").(string) != "" && d.Get("location").(string) != parts[2] {
			return fmt.Errorf("The location set on this secret version (%s) is not equal to the location where this secret exists (%s).", d.Get("location").(string), parts[2])
		}
		if err := d.Set("location", parts[2]); err != nil {
			return fmt.Errorf("Error setting location: %s", err)
		}
		if err := d.Set("secret", parts[3]); err != nil {
			return fmt.Errorf("Error setting secret: %s", err)
		}
	} else { // if secret name is provided in the secret field
		// Store values of project to set in state
		project, err = tpgresource.GetProject(d, config)
		if err != nil {
			return fmt.Errorf("Error fetching project for Secret: %s", err)
		}
		if d.Get("location").(string) == "" {
			return fmt.Errorf("Location must be set when providing only secret name")
		}
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}

	var url string
	versionNum := d.Get("version")

	// set version if provided, else set version to latest
	if versionNum != "" {
		url, err = tpgresource.ReplaceVars(d, config, "{{SecretManagerRegionalBasePath}}projects/{{project}}/locations/{{location}}/secrets/{{secret}}/versions/{{version}}")
		if err != nil {
			return err
		}
	} else {
		url, err = tpgresource.ReplaceVars(d, config, "{{SecretManagerRegionalBasePath}}projects/{{project}}/locations/{{location}}/secrets/{{secret}}/versions/latest")
		if err != nil {
			return err
		}
	}

	var secretVersion map[string]interface{}
	secretVersion, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   project,
		RawURL:    url,
		UserAgent: userAgent,
	})

	if err != nil {
		return fmt.Errorf("Error retrieving available secret manager regional secret versions: %s", err.Error())
	}

	secretVersionRegex := regexp.MustCompile("projects/(.+)/locations/(.+)/secrets/(.+)/versions/(.+)$")
	parts = secretVersionRegex.FindStringSubmatch(secretVersion["name"].(string))

	if len(parts) != 5 {
		return fmt.Errorf("secret name, %s, does not match format, projects/{{project}}/locations/{{location}}/secrets/{{secret}}/versions/{{version}}", secretVersion["name"].(string))
	}

	log.Printf("[DEBUG] Received Google Secret Manager Regional Secret Version: %q", secretVersion)

	if err := d.Set("version", parts[4]); err != nil {
		return fmt.Errorf("Error setting version: %s", err)
	}

	url = fmt.Sprintf("%s:access", url)
	resp, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   project,
		RawURL:    url,
		UserAgent: userAgent,
	})

	if err != nil {
		return fmt.Errorf("Error retrieving available secret manager regional secret version access: %s", err.Error())
	}

	if err := d.Set("customer_managed_encryption", flattenSecretManagerRegionalRegionalSecretVersionCustomerManagedEncryption(secretVersion["customerManagedEncryption"], d, config)); err != nil {
		return fmt.Errorf("Error setting customer_managed_encryption: %s", err)
	}

	if err := d.Set("create_time", secretVersion["createTime"].(string)); err != nil {
		return fmt.Errorf("Error setting create_time: %s", err)
	}

	if secretVersion["destroyTime"] != nil {
		if err := d.Set("destroy_time", secretVersion["destroyTime"].(string)); err != nil {
			return fmt.Errorf("Error setting destroy_time: %s", err)
		}
	}

	if err := d.Set("name", secretVersion["name"].(string)); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}

	if err := d.Set("enabled", true); err != nil {
		return fmt.Errorf("Error setting enabled: %s", err)
	}

	data := resp["payload"].(map[string]interface{})
	var secretData string
	if d.Get("is_secret_data_base64").(bool) {
		secretData = data["data"].(string)
	} else {
		payloadData, err := base64.StdEncoding.DecodeString(data["data"].(string))
		if err != nil {
			return fmt.Errorf("error decoding secret manager regional secret version data: %s", err.Error())
		}
		secretData = string(payloadData)
	}
	if err := d.Set("secret_data", secretData); err != nil {
		return fmt.Errorf("Error setting secret_data: %s", err)
	}

	d.SetId(secretVersion["name"].(string))
	return nil
}
