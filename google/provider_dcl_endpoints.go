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

package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// empty string is passed for dcl default since dcl
// [hardcodes the values](https://github.com/GoogleCloudPlatform/declarative-resource-client-library/blob/main/services/google/eventarc/beta/trigger_internal.go#L96-L103)

var AssuredWorkloadsEndpointEntryKey = "assured_workloads_custom_endpoint"
var AssuredWorkloadsEndpointEntry = &schema.Schema{
	Type:     schema.TypeString,
	Optional: true,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_ASSURED_WORKLOADS_CUSTOM_ENDPOINT",
	}, ""),
}

var ComputeEndpointEntryKey = "compute_custom_endpoint"
var ComputeEndpointEntry = &schema.Schema{
	Type:     schema.TypeString,
	Optional: true,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_COMPUTE_CUSTOM_ENDPOINT",
	}, ""),
}

var EventarcEndpointEntryKey = "eventarc_custom_endpoint"
var EventarcEndpointEntry = &schema.Schema{
	Type:     schema.TypeString,
	Optional: true,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_EVENTARC_CUSTOM_ENDPOINT",
	}, ""),
}

var PrivatecaEndpointEntryKey = "privateca_custom_endpoint"
var PrivatecaEndpointEntry = &schema.Schema{
	Type:     schema.TypeString,
	Optional: true,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_PRIVATECA_CUSTOM_ENDPOINT",
	}, ""),
}

//Add new values to config.go.erb config object declaration
//AssuredWorkloadsBasePath string
//ComputeBasePath string
//EventarcBasePath string
//PrivatecaBasePath string

//Add new values to provider.go.erb schema initialization
// AssuredWorkloadsEndpointEntryKey:               AssuredWorkloadsEndpointEntry,
// ComputeEndpointEntryKey:               ComputeEndpointEntry,
// EventarcEndpointEntryKey:               EventarcEndpointEntry,
// PrivatecaEndpointEntryKey:               PrivatecaEndpointEntry,

//Add new values to provider.go.erb - provider block read
// config.AssuredWorkloadsBasePath = d.Get(AssuredWorkloadsEndpointEntryKey).(string)
// config.ComputeBasePath = d.Get(ComputeEndpointEntryKey).(string)
// config.EventarcBasePath = d.Get(EventarcEndpointEntryKey).(string)
// config.PrivatecaBasePath = d.Get(PrivatecaEndpointEntryKey).(string)
