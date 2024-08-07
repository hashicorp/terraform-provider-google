// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: DCL     ***
//
// ----------------------------------------------------------------------------
//
//     This file is managed by Magic Modules (https://github.com/GoogleCloudPlatform/magic-modules)
//     and is based on the DCL (https://github.com/GoogleCloudPlatform/declarative-resource-client-library).
//     Changes will need to be made to the DCL or Magic Modules instead of here.
//
//     We are not currently able to accept contributions to this file. If changes
//     are required, please file an issue at https://github.com/hashicorp/terraform-provider-google/issues/new/choose
//
// ----------------------------------------------------------------------------

package transport

import (
	framework_schema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// empty string is passed for dcl default since dcl
// [hardcodes the values](https://github.com/GoogleCloudPlatform/declarative-resource-client-library/blob/main/services/google/eventarc/beta/trigger_internal.go#L96-L103)

var ApikeysEndpointEntryKey = "apikeys_custom_endpoint"
var ApikeysEndpointEntry = &schema.Schema{
	Type:     schema.TypeString,
	Optional: true,
}

var AssuredWorkloadsEndpointEntryKey = "assured_workloads_custom_endpoint"
var AssuredWorkloadsEndpointEntry = &schema.Schema{
	Type:     schema.TypeString,
	Optional: true,
}

var CloudBuildWorkerPoolEndpointEntryKey = "cloud_build_worker_pool_custom_endpoint"
var CloudBuildWorkerPoolEndpointEntry = &schema.Schema{
	Type:     schema.TypeString,
	Optional: true,
}

var CloudResourceManagerEndpointEntryKey = "cloud_resource_manager_custom_endpoint"
var CloudResourceManagerEndpointEntry = &schema.Schema{
	Type:     schema.TypeString,
	Optional: true,
}

var EventarcEndpointEntryKey = "eventarc_custom_endpoint"
var EventarcEndpointEntry = &schema.Schema{
	Type:     schema.TypeString,
	Optional: true,
}

var FirebaserulesEndpointEntryKey = "firebaserules_custom_endpoint"
var FirebaserulesEndpointEntry = &schema.Schema{
	Type:     schema.TypeString,
	Optional: true,
}

var GKEHubFeatureEndpointEntryKey = "gkehub_feature_custom_endpoint"
var GKEHubFeatureEndpointEntry = &schema.Schema{
	Type:     schema.TypeString,
	Optional: true,
}

var RecaptchaEnterpriseEndpointEntryKey = "recaptcha_enterprise_custom_endpoint"
var RecaptchaEnterpriseEndpointEntry = &schema.Schema{
	Type:     schema.TypeString,
	Optional: true,
}

type DCLConfig struct {
	ApikeysBasePath              string
	AssuredWorkloadsBasePath     string
	CloudBuildWorkerPoolBasePath string
	CloudResourceManagerBasePath string
	EventarcBasePath             string
	FirebaserulesBasePath        string
	GKEHubFeatureBasePath        string
	RecaptchaEnterpriseBasePath  string
}

func ConfigureDCLProvider(provider *schema.Provider) {
	provider.Schema[ApikeysEndpointEntryKey] = ApikeysEndpointEntry
	provider.Schema[AssuredWorkloadsEndpointEntryKey] = AssuredWorkloadsEndpointEntry
	provider.Schema[CloudBuildWorkerPoolEndpointEntryKey] = CloudBuildWorkerPoolEndpointEntry
	provider.Schema[CloudResourceManagerEndpointEntryKey] = CloudResourceManagerEndpointEntry
	provider.Schema[EventarcEndpointEntryKey] = EventarcEndpointEntry
	provider.Schema[FirebaserulesEndpointEntryKey] = FirebaserulesEndpointEntry
	provider.Schema[GKEHubFeatureEndpointEntryKey] = GKEHubFeatureEndpointEntry
	provider.Schema[RecaptchaEnterpriseEndpointEntryKey] = RecaptchaEnterpriseEndpointEntry
}

func HandleDCLCustomEndpointDefaults(d *schema.ResourceData) {
	if d.Get(ApikeysEndpointEntryKey) == "" {
		d.Set(ApikeysEndpointEntryKey, MultiEnvDefault([]string{
			"GOOGLE_APIKEYS_CUSTOM_ENDPOINT",
		}, ""))
	}
	if d.Get(AssuredWorkloadsEndpointEntryKey) == "" {
		d.Set(AssuredWorkloadsEndpointEntryKey, MultiEnvDefault([]string{
			"GOOGLE_ASSURED_WORKLOADS_CUSTOM_ENDPOINT",
		}, ""))
	}
	if d.Get(CloudBuildWorkerPoolEndpointEntryKey) == "" {
		d.Set(CloudBuildWorkerPoolEndpointEntryKey, MultiEnvDefault([]string{
			"GOOGLE_CLOUD_BUILD_WORKER_POOL_CUSTOM_ENDPOINT",
		}, ""))
	}
	if d.Get(CloudResourceManagerEndpointEntryKey) == "" {
		d.Set(CloudResourceManagerEndpointEntryKey, MultiEnvDefault([]string{
			"GOOGLE_CLOUD_RESOURCE_MANAGER_CUSTOM_ENDPOINT",
		}, ""))
	}
	if d.Get(EventarcEndpointEntryKey) == "" {
		d.Set(EventarcEndpointEntryKey, MultiEnvDefault([]string{
			"GOOGLE_EVENTARC_CUSTOM_ENDPOINT",
		}, ""))
	}
	if d.Get(FirebaserulesEndpointEntryKey) == "" {
		d.Set(FirebaserulesEndpointEntryKey, MultiEnvDefault([]string{
			"GOOGLE_FIREBASERULES_CUSTOM_ENDPOINT",
		}, ""))
	}
	if d.Get(GKEHubFeatureEndpointEntryKey) == "" {
		d.Set(GKEHubFeatureEndpointEntryKey, MultiEnvDefault([]string{
			"GOOGLE_GKEHUB_FEATURE_CUSTOM_ENDPOINT",
		}, ""))
	}
	if d.Get(RecaptchaEnterpriseEndpointEntryKey) == "" {
		d.Set(RecaptchaEnterpriseEndpointEntryKey, MultiEnvDefault([]string{
			"GOOGLE_RECAPTCHA_ENTERPRISE_CUSTOM_ENDPOINT",
		}, ""))
	}
}

// plugin-framework provider set-up
func ConfigureDCLCustomEndpointAttributesFramework(frameworkSchema *framework_schema.Schema) {
	frameworkSchema.Attributes["apikeys_custom_endpoint"] = framework_schema.StringAttribute{
		Optional: true,
		Validators: []validator.String{
			CustomEndpointValidator(),
		},
	}
	frameworkSchema.Attributes["assured_workloads_custom_endpoint"] = framework_schema.StringAttribute{
		Optional: true,
		Validators: []validator.String{
			CustomEndpointValidator(),
		},
	}
	frameworkSchema.Attributes["cloud_build_worker_pool_custom_endpoint"] = framework_schema.StringAttribute{
		Optional: true,
		Validators: []validator.String{
			CustomEndpointValidator(),
		},
	}
	frameworkSchema.Attributes["cloud_resource_manager_custom_endpoint"] = framework_schema.StringAttribute{
		Optional: true,
		Validators: []validator.String{
			CustomEndpointValidator(),
		},
	}
	frameworkSchema.Attributes["eventarc_custom_endpoint"] = framework_schema.StringAttribute{
		Optional: true,
		Validators: []validator.String{
			CustomEndpointValidator(),
		},
	}
	frameworkSchema.Attributes["firebaserules_custom_endpoint"] = framework_schema.StringAttribute{
		Optional: true,
		Validators: []validator.String{
			CustomEndpointValidator(),
		},
	}
	frameworkSchema.Attributes["gkehub_feature_custom_endpoint"] = framework_schema.StringAttribute{
		Optional: true,
		Validators: []validator.String{
			CustomEndpointValidator(),
		},
	}
	frameworkSchema.Attributes["recaptcha_enterprise_custom_endpoint"] = framework_schema.StringAttribute{
		Optional: true,
		Validators: []validator.String{
			CustomEndpointValidator(),
		},
	}
}

func ProviderDCLConfigure(d *schema.ResourceData, config *Config) interface{} {
	// networkConnectivity uses mmv1 basePath, assuredworkloads has a location variable in the basepath, can't be defined here.
	config.ApikeysBasePath = "https://apikeys.googleapis.com/v2/"
	config.AssuredWorkloadsBasePath = d.Get(AssuredWorkloadsEndpointEntryKey).(string)
	config.CloudBuildWorkerPoolBasePath = "https://cloudbuild.googleapis.com/v1/"
	config.CloudResourceManagerBasePath = "https://cloudresourcemanager.googleapis.com/"
	config.EventarcBasePath = "https://eventarc.googleapis.com/v1/"
	config.FirebaserulesBasePath = "https://firebaserules.googleapis.com/v1/"
	config.GKEHubFeatureBasePath = "https://gkehub.googleapis.com/v1beta1/"
	config.RecaptchaEnterpriseBasePath = "https://recaptchaenterprise.googleapis.com/v1/"

	return config
}
