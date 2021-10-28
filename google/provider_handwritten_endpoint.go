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

var DataprocBetaCustomEndpointEntryKey = "dataproc_beta_custom_endpoint"
var DataprocBetaCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: validateCustomEndpoint,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_DATAPROC_BETA_CUSTOM_ENDPOINT",
	}, DefaultBasePaths[DataprocBetaBasePathKey]),
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

var ResourceManagerV2CustomEndpointEntryKey = "resource_manager_v2_custom_endpoint"
var ResourceManagerV2CustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: validateCustomEndpoint,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_RESOURCE_MANAGER_V2_CUSTOM_ENDPOINT",
	}, DefaultBasePaths[ResourceManagerV2BasePathKey]),
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

var StorageTransferCustomEndpointEntryKey = "storage_transfer_custom_endpoint"
var StorageTransferCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: validateCustomEndpoint,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_STORAGE_TRANSFER_CUSTOM_ENDPOINT",
	}, DefaultBasePaths[StorageTransferBasePathKey]),
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

// GkeHubFeature uses a different base path "v1beta" than GkeHubMembership "v1beta1"
var GkeHubFeatureCustomEndpointEntryKey = "gkehub_feature_custom_endpoint"
var GkeHubFeatureCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: validateCustomEndpoint,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_GKEHUB_FEATURE_CUSTOM_ENDPOINT",
	}, DefaultBasePaths[GkeHubFeatureBasePathKey]),
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

func validateCustomEndpoint(v interface{}, k string) (ws []string, errors []error) {
	re := `.*/[^/]+/$`
	return validateRegexp(re)(v, k)
}
