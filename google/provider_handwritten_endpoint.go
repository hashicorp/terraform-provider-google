package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// For generated resources, endpoint entries live in product-specific provider
// files. Collect handwritten ones here. If any of these are modified, be sure
// to update the provider_reference docs page.

var CloudBillingDefaultBasePath = "https://cloudbilling.googleapis.com/v1/"
var CloudBillingCustomEndpointEntryKey = "cloud_billing_custom_endpoint"
var CloudBillingCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: validateCustomEndpoint,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_CLOUD_BILLING_CUSTOM_ENDPOINT",
	}, CloudBillingDefaultBasePath),
}

var CloudIoTDefaultBasePath = "https://cloudiot.googleapis.com/v1/"
var CloudIoTCustomEndpointEntryKey = "cloud_iot_custom_endpoint"
var CloudIoTCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: validateCustomEndpoint,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_CLOUD_IOT_CUSTOM_ENDPOINT",
	}, CloudIoTDefaultBasePath),
}

var ComposerDefaultBasePath = "https://composer.googleapis.com/v1beta1/"
var ComposerCustomEndpointEntryKey = "composer_custom_endpoint"
var ComposerCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: validateCustomEndpoint,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_COMPOSER_CUSTOM_ENDPOINT",
	}, ComposerDefaultBasePath),
}

var ComputeBetaDefaultBasePath = "https://www.googleapis.com/compute/beta/"
var ComputeBetaCustomEndpointEntryKey = "compute_beta_custom_endpoint"
var ComputeBetaCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: validateCustomEndpoint,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_COMPUTE_BETA_CUSTOM_ENDPOINT",
	}, ComputeBetaDefaultBasePath),
}

var ContainerDefaultBasePath = "https://container.googleapis.com/v1/"
var ContainerCustomEndpointEntryKey = "container_custom_endpoint"
var ContainerCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: validateCustomEndpoint,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_CONTAINER_CUSTOM_ENDPOINT",
	}, ContainerDefaultBasePath),
}

var ContainerBetaDefaultBasePath = "https://container.googleapis.com/v1beta1/"
var ContainerBetaCustomEndpointEntryKey = "container_beta_custom_endpoint"
var ContainerBetaCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: validateCustomEndpoint,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_CONTAINER_BETA_CUSTOM_ENDPOINT",
	}, ContainerBetaDefaultBasePath),
}

var DataprocBetaDefaultBasePath = "https://dataproc.googleapis.com/v1beta2/"
var DataprocBetaCustomEndpointEntryKey = "dataproc_beta_custom_endpoint"
var DataprocBetaCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: validateCustomEndpoint,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_DATAPROC_BETA_CUSTOM_ENDPOINT",
	}, DataprocBetaDefaultBasePath),
}

var DataflowDefaultBasePath = "https://dataflow.googleapis.com/v1b3/"
var DataflowCustomEndpointEntryKey = "dataflow_custom_endpoint"
var DataflowCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: validateCustomEndpoint,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_DATAFLOW_CUSTOM_ENDPOINT",
	}, DataflowDefaultBasePath),
}
var DnsBetaDefaultBasePath = "https://www.googleapis.com/dns/v1beta2/"
var DnsBetaCustomEndpointEntryKey = "dns_beta_custom_endpoint"
var DnsBetaCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: validateCustomEndpoint,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_DNS_BETA_CUSTOM_ENDPOINT",
	}, DnsBetaDefaultBasePath),
}

var IAMDefaultBasePath = "https://iam.googleapis.com/v1/"
var IAMCustomEndpointEntryKey = "iam_custom_endpoint"
var IAMCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: validateCustomEndpoint,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_IAM_CUSTOM_ENDPOINT",
	}, IAMDefaultBasePath),
}

var IamCredentialsDefaultBasePath = "https://iamcredentials.googleapis.com/v1/"
var IamCredentialsCustomEndpointEntryKey = "iam_credentials_custom_endpoint"
var IamCredentialsCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: validateCustomEndpoint,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_IAM_CREDENTIALS_CUSTOM_ENDPOINT",
	}, IamCredentialsDefaultBasePath),
}

var ResourceManagerV2Beta1DefaultBasePath = "https://cloudresourcemanager.googleapis.com/v2beta1/"
var ResourceManagerV2Beta1CustomEndpointEntryKey = "resource_manager_v2beta1_custom_endpoint"
var ResourceManagerV2Beta1CustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: validateCustomEndpoint,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_RESOURCE_MANAGER_V2BETA1_CUSTOM_ENDPOINT",
	}, ResourceManagerV2Beta1DefaultBasePath),
}

var RuntimeConfigCustomEndpointEntryKey = "runtimeconfig_custom_endpoint"
var RuntimeConfigCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: validateCustomEndpoint,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_RUNTIMECONFIG_CUSTOM_ENDPOINT",
	}, RuntimeConfigDefaultBasePath),
}

var ServiceManagementDefaultBasePath = "https://servicemanagement.googleapis.com/v1/"
var ServiceManagementCustomEndpointEntryKey = "service_management_custom_endpoint"
var ServiceManagementCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: validateCustomEndpoint,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_SERVICE_MANAGEMENT_CUSTOM_ENDPOINT",
	}, ServiceManagementDefaultBasePath),
}

var ServiceNetworkingDefaultBasePath = "https://servicenetworking.googleapis.com/v1/"
var ServiceNetworkingCustomEndpointEntryKey = "service_networking_custom_endpoint"
var ServiceNetworkingCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: validateCustomEndpoint,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_SERVICE_NETWORKING_CUSTOM_ENDPOINT",
	}, ServiceNetworkingDefaultBasePath),
}

var ServiceUsageDefaultBasePath = "https://serviceusage.googleapis.com/v1/"
var ServiceUsageCustomEndpointEntryKey = "service_usage_custom_endpoint"
var ServiceUsageCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: validateCustomEndpoint,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_SERVICE_USAGE_CUSTOM_ENDPOINT",
	}, ServiceUsageDefaultBasePath),
}

var StorageTransferDefaultBasePath = "https://storagetransfer.googleapis.com/v1/"
var StorageTransferCustomEndpointEntryKey = "storage_transfer_custom_endpoint"
var StorageTransferCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: validateCustomEndpoint,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_STORAGE_TRANSFER_CUSTOM_ENDPOINT",
	}, StorageTransferDefaultBasePath),
}

var BigtableAdminDefaultBasePath = "https://bigtableadmin.googleapis.com/v2/"
var BigtableAdminCustomEndpointEntryKey = "bigtable_custom_endpoint"
var BigtableAdminCustomEndpointEntry = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: validateCustomEndpoint,
	DefaultFunc: schema.MultiEnvDefaultFunc([]string{
		"GOOGLE_BIGTABLE_CUSTOM_ENDPOINT",
	}, BigtableAdminDefaultBasePath),
}

func validateCustomEndpoint(v interface{}, k string) (ws []string, errors []error) {
	re := `.*/[^/]+/$`
	return validateRegexp(re)(v, k)
}
