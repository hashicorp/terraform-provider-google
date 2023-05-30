// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package transport

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/verify"
)

// For generated resources, endpoint entries live in product-specific provider
// files. Collect handwritten ones here. If any of these are modified, be sure
// to update the provider_reference docs page.

var CloudBillingCustomEndpointEntryKey = "cloud_billing_custom_endpoint"
var CloudBillingCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: ValidateCustomEndpoint,
}

var ComposerCustomEndpointEntryKey = "composer_custom_endpoint"
var ComposerCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: ValidateCustomEndpoint,
}

var ContainerCustomEndpointEntryKey = "container_custom_endpoint"
var ContainerCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: ValidateCustomEndpoint,
}

var DataflowCustomEndpointEntryKey = "dataflow_custom_endpoint"
var DataflowCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: ValidateCustomEndpoint,
}

var IAMCustomEndpointEntryKey = "iam_custom_endpoint"
var IAMCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: ValidateCustomEndpoint,
}

var IamCredentialsCustomEndpointEntryKey = "iam_credentials_custom_endpoint"
var IamCredentialsCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: ValidateCustomEndpoint,
}

var ResourceManagerV3CustomEndpointEntryKey = "resource_manager_v3_custom_endpoint"
var ResourceManagerV3CustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: ValidateCustomEndpoint,
}

var ServiceNetworkingCustomEndpointEntryKey = "service_networking_custom_endpoint"
var ServiceNetworkingCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: ValidateCustomEndpoint,
}

var ServiceUsageCustomEndpointEntryKey = "service_usage_custom_endpoint"
var ServiceUsageCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: ValidateCustomEndpoint,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_SERVICE_USAGE_CUSTOM_ENDPOINT",
	}, DefaultBasePaths[ServiceUsageBasePathKey]),
}

var BigtableAdminCustomEndpointEntryKey = "bigtable_custom_endpoint"
var BigtableAdminCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: ValidateCustomEndpoint,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_BIGTABLE_CUSTOM_ENDPOINT",
	}, DefaultBasePaths[BigtableAdminBasePathKey]),
}

var PrivatecaCertificateTemplateEndpointEntryKey = "privateca_custom_endpoint"
var PrivatecaCertificateTemplateCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: ValidateCustomEndpoint,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_PRIVATECA_CUSTOM_ENDPOINT",
	}, DefaultBasePaths[PrivatecaBasePathKey]),
}

var ContainerAwsCustomEndpointEntryKey = "container_aws_custom_endpoint"
var ContainerAwsCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: ValidateCustomEndpoint,
}

var ContainerAzureCustomEndpointEntryKey = "container_azure_custom_endpoint"
var ContainerAzureCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: ValidateCustomEndpoint,
}

var TagsLocationCustomEndpointEntryKey = "tags_location_custom_endpoint"
var TagsLocationCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: ValidateCustomEndpoint,
}

func ValidateCustomEndpoint(v interface{}, k string) (ws []string, errors []error) {
	re := `.*/[^/]+/$`
	return verify.ValidateRegexp(re)(v, k)
}
