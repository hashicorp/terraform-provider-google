package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// For generated resources, endpoint entries live in product-specific provider
// files. Collect handwritten ones here. If any of these are modified, be sure
// to update the provider_reference docs page.

var CloudBillingCustomEndpointEntryKey = "cloud_billing_custom_endpoint"
var CloudBillingCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: validateCustomEndpoint,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_CLOUD_BILLING_CUSTOM_ENDPOINT",
	}, DefaultBasePaths[CloudBillingBasePathKey]),
}

var ComposerCustomEndpointEntryKey = "composer_custom_endpoint"
var ComposerCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: validateCustomEndpoint,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_COMPOSER_CUSTOM_ENDPOINT",
	}, DefaultBasePaths[ComposerBasePathKey]),
}

var ContainerCustomEndpointEntryKey = "container_custom_endpoint"
var ContainerCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: validateCustomEndpoint,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_CONTAINER_CUSTOM_ENDPOINT",
	}, DefaultBasePaths[ContainerBasePathKey]),
}

var DataflowCustomEndpointEntryKey = "dataflow_custom_endpoint"
var DataflowCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: validateCustomEndpoint,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_DATAFLOW_CUSTOM_ENDPOINT",
	}, DefaultBasePaths[DataflowBasePathKey]),
}

var IAMCustomEndpointEntryKey = "iam_custom_endpoint"
var IAMCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: validateCustomEndpoint,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_IAM_CUSTOM_ENDPOINT",
	}, DefaultBasePaths[IAMBasePathKey]),
}

var IamCredentialsCustomEndpointEntryKey = "iam_credentials_custom_endpoint"
var IamCredentialsCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: validateCustomEndpoint,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_IAM_CREDENTIALS_CUSTOM_ENDPOINT",
	}, DefaultBasePaths[IamCredentialsBasePathKey]),
}

var ResourceManagerV3CustomEndpointEntryKey = "resource_manager_v3_custom_endpoint"
var ResourceManagerV3CustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: validateCustomEndpoint,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_RESOURCE_MANAGER_V3_CUSTOM_ENDPOINT",
	}, DefaultBasePaths[ResourceManagerV3BasePathKey]),
}

var ServiceNetworkingCustomEndpointEntryKey = "service_networking_custom_endpoint"
var ServiceNetworkingCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: validateCustomEndpoint,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_SERVICE_NETWORKING_CUSTOM_ENDPOINT",
	}, DefaultBasePaths[ServiceNetworkingBasePathKey]),
}

var ServiceUsageCustomEndpointEntryKey = "service_usage_custom_endpoint"
var ServiceUsageCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: validateCustomEndpoint,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_SERVICE_USAGE_CUSTOM_ENDPOINT",
	}, DefaultBasePaths[ServiceUsageBasePathKey]),
}

var BigtableAdminCustomEndpointEntryKey = "bigtable_custom_endpoint"
var BigtableAdminCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: validateCustomEndpoint,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_BIGTABLE_CUSTOM_ENDPOINT",
	}, DefaultBasePaths[BigtableAdminBasePathKey]),
}

var PrivatecaCertificateTemplateEndpointEntryKey = "privateca_custom_endpoint"
var PrivatecaCertificateTemplateCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: validateCustomEndpoint,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_PRIVATECA_CUSTOM_ENDPOINT",
	}, DefaultBasePaths[PrivatecaBasePathKey]),
}

var ContainerAwsCustomEndpointEntryKey = "container_aws_custom_endpoint"
var ContainerAwsCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: validateCustomEndpoint,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_CONTAINERAWS_CUSTOM_ENDPOINT",
	}, DefaultBasePaths[ContainerAwsBasePathKey]),
}

var ContainerAzureCustomEndpointEntryKey = "container_azure_custom_endpoint"
var ContainerAzureCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: validateCustomEndpoint,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_CONTAINERAZURE_CUSTOM_ENDPOINT",
	}, DefaultBasePaths[ContainerAzureBasePathKey]),
}

func validateCustomEndpoint(v interface{}, k string) (ws []string, errors []error) {
	re := `.*/[^/]+/$`
	return validateRegexp(re)(v, k)
}
