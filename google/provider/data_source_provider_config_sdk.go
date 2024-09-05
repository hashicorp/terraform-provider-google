// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package provider

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleProviderConfigSdk() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "Data source google_provider_config_sdk is intended to be used only in acceptance tests for the provider. Instead, please use the google_client_config data source to access provider configuration details, or open a GitHub issue requesting new features in that datasource. Please go to: https://github.com/hashicorp/terraform-provider-google/issues/new/choose",
		Read:               dataSourceClientConfigRead,
		Schema: map[string]*schema.Schema{
			// Start of user inputs
			"access_token": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"credentials": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"impersonate_service_account": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"impersonate_service_account_delegates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"project": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"billing_project": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"universe_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scopes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"user_project_override": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"request_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_timeout": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_labels": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"add_terraform_attribution_label": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"terraform_attribution_label_addition_strategy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// End of user inputs

			// Note - this data source excludes the default and custom endpoints for individual services

			// Start of values set during provider configuration
			"user_agent": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// End of values set during provider configuration
		},
	}
}

func dataSourceClientConfigRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	if err := d.Set("access_token", config.AccessToken); err != nil {
		return fmt.Errorf("error setting access_token: %s", err)
	}
	if err := d.Set("credentials", config.Credentials); err != nil {
		return fmt.Errorf("error setting credentials: %s", err)
	}
	if err := d.Set("impersonate_service_account", config.ImpersonateServiceAccount); err != nil {
		return fmt.Errorf("error setting impersonate_service_account: %s", err)
	}
	if err := d.Set("impersonate_service_account_delegates", config.ImpersonateServiceAccountDelegates); err != nil {
		return fmt.Errorf("error setting impersonate_service_account_delegates: %s", err)
	}
	if err := d.Set("project", config.Project); err != nil {
		return fmt.Errorf("error setting project: %s", err)
	}
	if err := d.Set("region", config.Region); err != nil {
		return fmt.Errorf("error setting region: %s", err)
	}
	if err := d.Set("billing_project", config.BillingProject); err != nil {
		return fmt.Errorf("error setting billing_project: %s", err)
	}
	if err := d.Set("zone", config.Zone); err != nil {
		return fmt.Errorf("error setting zone: %s", err)
	}
	if err := d.Set("universe_domain", config.UniverseDomain); err != nil {
		return fmt.Errorf("error setting universe_domain: %s", err)
	}
	if err := d.Set("scopes", config.Scopes); err != nil {
		return fmt.Errorf("error setting scopes: %s", err)
	}
	if err := d.Set("user_project_override", config.UserProjectOverride); err != nil {
		return fmt.Errorf("error setting user_project_override: %s", err)
	}
	if err := d.Set("request_reason", config.RequestReason); err != nil {
		return fmt.Errorf("error setting request_reason: %s", err)
	}
	if err := d.Set("request_timeout", config.RequestTimeout.String()); err != nil {
		return fmt.Errorf("error setting request_timeout: %s", err)
	}
	if err := d.Set("default_labels", config.DefaultLabels); err != nil {
		return fmt.Errorf("error setting default_labels: %s", err)
	}
	if err := d.Set("add_terraform_attribution_label", config.AddTerraformAttributionLabel); err != nil {
		return fmt.Errorf("error setting add_terraform_attribution_label: %s", err)
	}
	if err := d.Set("terraform_attribution_label_addition_strategy", config.TerraformAttributionLabelAdditionStrategy); err != nil {
		return fmt.Errorf("error setting terraform_attribution_label_addition_strategy: %s", err)
	}
	if err := d.Set("user_agent", config.UserAgent); err != nil {
		return fmt.Errorf("error setting user_agent: %s", err)
	}

	// Id is a hash of the total transport.Config struct
	configString := []byte(fmt.Sprintf("%#v", config))
	hasher := sha1.New()
	hasher.Write(configString)
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	d.SetId(string(sha))

	return nil
}
