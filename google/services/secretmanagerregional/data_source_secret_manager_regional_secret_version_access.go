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

func DataSourceSecretManagerRegionalRegionalSecretVersionAccess() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSecretManagerRegionalRegionalSecretVersionAccessRead,
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
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret_data": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}
func dataSourceSecretManagerRegionalRegionalSecretVersionAccessRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	secretRegex := regexp.MustCompile("projects/(.+)/locations/(.+)/secrets/(.+)$")
	dSecret, ok := d.Get("secret").(string)
	if !ok {
		return fmt.Errorf("wrong type for secret field (%T), expected string", d.Get("secret"))
	}
	parts := secretRegex.FindStringSubmatch(dSecret)

	var project string

	// if reference of the secret is provided in the secret field
	if len(parts) == 4 {
		// Stores value of project to set in state
		project = parts[1]
		if dProject, ok := d.Get("project").(string); !ok {
			return fmt.Errorf("wrong type for project (%T), expected string", d.Get("project"))
		} else if dProject != "" && dProject != project {
			return fmt.Errorf("project field value (%s) does not match project of secret (%s).", d.Get("project").(string), project)
		}
		if dLocation, ok := d.Get("location").(string); !ok {
			return fmt.Errorf("wrong type for location (%T), expected string", d.Get("location"))
		} else if dLocation != "" && dLocation != parts[2] {
			return fmt.Errorf("location field value (%s) does not match location of secret (%s).", dLocation, parts[2])
		}
		if err := d.Set("location", parts[2]); err != nil {
			return fmt.Errorf("error setting location: %s", err)
		}
		if err := d.Set("secret", parts[3]); err != nil {
			return fmt.Errorf("error setting secret: %s", err)
		}
	} else { // if secret name is provided in the secret field
		// Stores value of project to set in state
		project, err = tpgresource.GetProject(d, config)
		if err != nil {
			return fmt.Errorf("error fetching project for Secret: %s", err)
		}
		if dLocation, ok := d.Get("location").(string); ok && dLocation == "" {
			return fmt.Errorf("location must be set when providing only secret name")
		}
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("error setting project: %s", err)
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

	url = fmt.Sprintf("%s:access", url)
	resp, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   project,
		RawURL:    url,
		UserAgent: userAgent,
	})

	if err != nil {
		return fmt.Errorf("error retrieving available secret manager regional secret version access: %s", err.Error())
	}

	nameValue, ok := resp["name"]
	if !ok {
		return fmt.Errorf("read response didn't contain critical fields. Read may not have succeeded.")
	}
	if err := d.Set("name", nameValue.(string)); err != nil {
		return fmt.Errorf("error setting name: %s", err)
	}

	secretVersionRegex := regexp.MustCompile("projects/(.+)/locations/(.+)/secrets/(.+)/versions/(.+)$")
	parts = secretVersionRegex.FindStringSubmatch(nameValue.(string))
	if len(parts) != 5 {
		return fmt.Errorf("secret name, %s, does not match format, projects/{{project}}/locations/{{location}}/secrets/{{secret}}/versions/{{version}}", nameValue.(string))
	}

	log.Printf("[DEBUG] Received Google SecretManager Version: %q", parts[3])

	if err := d.Set("version", parts[4]); err != nil {
		return fmt.Errorf("error setting version: %s", err)
	}

	data := resp["payload"].(map[string]interface{})
	secretData, err := base64.StdEncoding.DecodeString(data["data"].(string))
	if err != nil {
		return fmt.Errorf("error decoding secret manager regional secret version data: %s", err.Error())
	}
	if err := d.Set("secret_data", string(secretData)); err != nil {
		return fmt.Errorf("error setting secret_data: %s", err)
	}

	d.SetId(nameValue.(string))
	return nil
}
