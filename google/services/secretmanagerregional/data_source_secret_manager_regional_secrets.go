// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package secretmanagerregional

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceSecretManagerRegionalRegionalSecrets() *schema.Resource {
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceSecretManagerRegionalRegionalSecret().Schema)
	return &schema.Resource{
		Read: dataSourceSecretManagerRegionalRegionalSecretsRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"location": {
				Type:     schema.TypeString,
				Required: true,
			},
			"filter": {
				Type: schema.TypeString,
				Description: `Filter string, adhering to the rules in List-operation filtering (https://cloud.google.com/secret-manager/docs/filtering).
List only secrets matching the filter. If filter is empty, all regional secrets are listed from the specified location.`,
				Optional: true,
			},
			"secrets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: dsSchema,
				},
			},
		},
	}
}

func dataSourceSecretManagerRegionalRegionalSecretsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{SecretManagerRegionalBasePath}}projects/{{project}}/locations/{{location}}/secrets")
	if err != nil {
		return err
	}

	filter, has_filter := d.GetOk("filter")

	if has_filter {
		url, err = transport_tpg.AddQueryParams(url, map[string]string{"filter": filter.(string)})
		if err != nil {
			return err
		}
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Secret: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	// To handle pagination locally
	allSecrets := make([]interface{}, 0)
	token := ""
	for paginate := true; paginate; {
		if token != "" {
			url, err = transport_tpg.AddQueryParams(url, map[string]string{"pageToken": token})
			if err != nil {
				return err
			}
		}
		secrets, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			Project:   billingProject,
			RawURL:    url,
			UserAgent: userAgent,
		})
		if err != nil {
			return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("SecretManagerRegionalSecrets %q", d.Id()))
		}
		secretsInterface := secrets["secrets"]
		if secretsInterface == nil {
			break
		}
		allSecrets = append(allSecrets, secretsInterface.([]interface{})...)
		tokenInterface := secrets["nextPageToken"]
		if tokenInterface == nil {
			paginate = false
		} else {
			paginate = true
			token = tokenInterface.(string)
		}
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("error setting project: %s", err)
	}

	if err := d.Set("filter", filter); err != nil {
		return fmt.Errorf("error setting filter: %s", err)
	}

	if err := d.Set("secrets", flattenSecretManagerRegionalRegionalSecretsSecrets(allSecrets, d, config)); err != nil {
		return fmt.Errorf("error setting secrets: %s", err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/secrets")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	if has_filter {
		id += "/filter=" + filter.(string)
	}
	d.SetId(id)

	return nil
}

func flattenSecretManagerRegionalRegionalSecretsSecrets(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))

	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Do not include empty json objects coming back from the api
			continue
		}

		transformed = append(transformed, map[string]interface{}{
			"annotations":                 flattenSecretManagerRegionalRegionalSecretEffectiveAnnotations(original["annotations"], d, config),
			"effective_annotations":       flattenSecretManagerRegionalRegionalSecretEffectiveAnnotations(original["annotations"], d, config),
			"expire_time":                 flattenSecretManagerRegionalRegionalSecretExpireTime(original["expireTime"], d, config),
			"labels":                      flattenSecretManagerRegionalRegionalSecretEffectiveLabels(original["labels"], d, config),
			"effective_labels":            flattenSecretManagerRegionalRegionalSecretEffectiveLabels(original["labels"], d, config),
			"terraform_labels":            flattenSecretManagerRegionalRegionalSecretEffectiveLabels(original["labels"], d, config),
			"version_aliases":             flattenSecretManagerRegionalRegionalSecretVersionAliases(original["versionAliases"], d, config),
			"rotation":                    flattenSecretManagerRegionalRegionalSecretRotation(original["rotation"], d, config),
			"topics":                      flattenSecretManagerRegionalRegionalSecretTopics(original["topics"], d, config),
			"version_destroy_ttl":         flattenSecretManagerRegionalRegionalSecretVersionDestroyTtl(original["versionDestroyTtl"], d, config),
			"customer_managed_encryption": flattenSecretManagerRegionalRegionalSecretCustomerManagedEncryption(original["customerManagedEncryption"], d, config),
			"create_time":                 flattenSecretManagerRegionalRegionalSecretCreateTime(original["createTime"], d, config),
			"name":                        flattenSecretManagerRegionalRegionalSecretName(original["name"], d, config),
			"project":                     getDataFromName(original["name"], 1),
			"location":                    getDataFromName(original["name"], 3),
			"secret_id":                   getDataFromName(original["name"], 5),
		})
	}
	return transformed
}

func getDataFromName(v interface{}, part int) string {
	name := v.(string)
	split := strings.Split(name, "/")
	return split[part]
}
