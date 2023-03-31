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
	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	"time"

	apikeys "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/apikeys"
	assuredworkloads "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/assuredworkloads"
	bigqueryreservation "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/bigqueryreservation"
	cloudbuild "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/cloudbuild"
	clouddeploy "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/clouddeploy"
	cloudresourcemanager "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/cloudresourcemanager"
	compute "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/compute"
	containeraws "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/containeraws"
	containerazure "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/containerazure"
	dataplex "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/dataplex"
	dataproc "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/dataproc"
	eventarc "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/eventarc"
	firebaserules "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/firebaserules"
	monitoring "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/monitoring"
	networkconnectivity "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/networkconnectivity"
	orgpolicy "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/orgpolicy"
	osconfig "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/osconfig"
	privateca "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/privateca"
	recaptchaenterprise "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/recaptchaenterprise"
)

func NewDCLApikeysClient(config *Config, userAgent, billingProject string, timeout time.Duration) *apikeys.Client {
	configOptions := []dcl.ConfigOption{
		dcl.WithHTTPClient(config.Client),
		dcl.WithUserAgent(userAgent),
		dcl.WithLogger(dclLogger{}),
		dcl.WithBasePath(config.ApikeysBasePath),
	}

	if timeout != 0 {
		configOptions = append(configOptions, dcl.WithTimeout(timeout))
	}

	if config.UserProjectOverride {
		configOptions = append(configOptions, dcl.WithUserProjectOverride())
		if billingProject != "" {
			configOptions = append(configOptions, dcl.WithBillingProject(billingProject))
		}
	}

	dclConfig := dcl.NewConfig(configOptions...)
	return apikeys.NewClient(dclConfig)
}

func NewDCLAssuredWorkloadsClient(config *Config, userAgent, billingProject string, timeout time.Duration) *assuredworkloads.Client {
	configOptions := []dcl.ConfigOption{
		dcl.WithHTTPClient(config.Client),
		dcl.WithUserAgent(userAgent),
		dcl.WithLogger(dclLogger{}),
		dcl.WithBasePath(config.AssuredWorkloadsBasePath),
	}

	if timeout != 0 {
		configOptions = append(configOptions, dcl.WithTimeout(timeout))
	}

	if config.UserProjectOverride {
		configOptions = append(configOptions, dcl.WithUserProjectOverride())
		if billingProject != "" {
			configOptions = append(configOptions, dcl.WithBillingProject(billingProject))
		}
	}

	dclConfig := dcl.NewConfig(configOptions...)
	return assuredworkloads.NewClient(dclConfig)
}

func NewDCLBigqueryReservationClient(config *Config, userAgent, billingProject string, timeout time.Duration) *bigqueryreservation.Client {
	configOptions := []dcl.ConfigOption{
		dcl.WithHTTPClient(config.Client),
		dcl.WithUserAgent(userAgent),
		dcl.WithLogger(dclLogger{}),
		dcl.WithBasePath(config.BigqueryReservationBasePath),
	}

	if timeout != 0 {
		configOptions = append(configOptions, dcl.WithTimeout(timeout))
	}

	if config.UserProjectOverride {
		configOptions = append(configOptions, dcl.WithUserProjectOverride())
		if billingProject != "" {
			configOptions = append(configOptions, dcl.WithBillingProject(billingProject))
		}
	}

	dclConfig := dcl.NewConfig(configOptions...)
	return bigqueryreservation.NewClient(dclConfig)
}

func NewDCLCloudbuildClient(config *Config, userAgent, billingProject string, timeout time.Duration) *cloudbuild.Client {
	configOptions := []dcl.ConfigOption{
		dcl.WithHTTPClient(config.Client),
		dcl.WithUserAgent(userAgent),
		dcl.WithLogger(dclLogger{}),
		dcl.WithBasePath(config.CloudBuildWorkerPoolBasePath),
	}

	if timeout != 0 {
		configOptions = append(configOptions, dcl.WithTimeout(timeout))
	}

	if config.UserProjectOverride {
		configOptions = append(configOptions, dcl.WithUserProjectOverride())
		if billingProject != "" {
			configOptions = append(configOptions, dcl.WithBillingProject(billingProject))
		}
	}

	dclConfig := dcl.NewConfig(configOptions...)
	return cloudbuild.NewClient(dclConfig)
}

func NewDCLClouddeployClient(config *Config, userAgent, billingProject string, timeout time.Duration) *clouddeploy.Client {
	configOptions := []dcl.ConfigOption{
		dcl.WithHTTPClient(config.Client),
		dcl.WithUserAgent(userAgent),
		dcl.WithLogger(dclLogger{}),
		dcl.WithBasePath(config.ClouddeployBasePath),
	}

	if timeout != 0 {
		configOptions = append(configOptions, dcl.WithTimeout(timeout))
	}

	if config.UserProjectOverride {
		configOptions = append(configOptions, dcl.WithUserProjectOverride())
		if billingProject != "" {
			configOptions = append(configOptions, dcl.WithBillingProject(billingProject))
		}
	}

	dclConfig := dcl.NewConfig(configOptions...)
	return clouddeploy.NewClient(dclConfig)
}

func NewDCLCloudResourceManagerClient(config *Config, userAgent, billingProject string, timeout time.Duration) *cloudresourcemanager.Client {
	configOptions := []dcl.ConfigOption{
		dcl.WithHTTPClient(config.Client),
		dcl.WithUserAgent(userAgent),
		dcl.WithLogger(dclLogger{}),
		dcl.WithBasePath(config.CloudResourceManagerBasePath),
	}

	if timeout != 0 {
		configOptions = append(configOptions, dcl.WithTimeout(timeout))
	}

	if config.UserProjectOverride {
		configOptions = append(configOptions, dcl.WithUserProjectOverride())
		if billingProject != "" {
			configOptions = append(configOptions, dcl.WithBillingProject(billingProject))
		}
	}

	dclConfig := dcl.NewConfig(configOptions...)
	return cloudresourcemanager.NewClient(dclConfig)
}

func NewDCLComputeClient(config *Config, userAgent, billingProject string, timeout time.Duration) *compute.Client {
	configOptions := []dcl.ConfigOption{
		dcl.WithHTTPClient(config.Client),
		dcl.WithUserAgent(userAgent),
		dcl.WithLogger(dclLogger{}),
		dcl.WithBasePath(config.ComputeBasePath),
	}

	if timeout != 0 {
		configOptions = append(configOptions, dcl.WithTimeout(timeout))
	}

	if config.UserProjectOverride {
		configOptions = append(configOptions, dcl.WithUserProjectOverride())
		if billingProject != "" {
			configOptions = append(configOptions, dcl.WithBillingProject(billingProject))
		}
	}

	dclConfig := dcl.NewConfig(configOptions...)
	return compute.NewClient(dclConfig)
}

func NewDCLContainerAwsClient(config *Config, userAgent, billingProject string, timeout time.Duration) *containeraws.Client {
	configOptions := []dcl.ConfigOption{
		dcl.WithHTTPClient(config.Client),
		dcl.WithUserAgent(userAgent),
		dcl.WithLogger(dclLogger{}),
		dcl.WithBasePath(config.ContainerAwsBasePath),
	}

	if timeout != 0 {
		configOptions = append(configOptions, dcl.WithTimeout(timeout))
	}

	if config.UserProjectOverride {
		configOptions = append(configOptions, dcl.WithUserProjectOverride())
		if billingProject != "" {
			configOptions = append(configOptions, dcl.WithBillingProject(billingProject))
		}
	}

	dclConfig := dcl.NewConfig(configOptions...)
	return containeraws.NewClient(dclConfig)
}

func NewDCLContainerAzureClient(config *Config, userAgent, billingProject string, timeout time.Duration) *containerazure.Client {
	configOptions := []dcl.ConfigOption{
		dcl.WithHTTPClient(config.Client),
		dcl.WithUserAgent(userAgent),
		dcl.WithLogger(dclLogger{}),
		dcl.WithBasePath(config.ContainerAzureBasePath),
	}

	if timeout != 0 {
		configOptions = append(configOptions, dcl.WithTimeout(timeout))
	}

	if config.UserProjectOverride {
		configOptions = append(configOptions, dcl.WithUserProjectOverride())
		if billingProject != "" {
			configOptions = append(configOptions, dcl.WithBillingProject(billingProject))
		}
	}

	dclConfig := dcl.NewConfig(configOptions...)
	return containerazure.NewClient(dclConfig)
}

func NewDCLDataplexClient(config *Config, userAgent, billingProject string, timeout time.Duration) *dataplex.Client {
	configOptions := []dcl.ConfigOption{
		dcl.WithHTTPClient(config.Client),
		dcl.WithUserAgent(userAgent),
		dcl.WithLogger(dclLogger{}),
		dcl.WithBasePath(config.DataplexBasePath),
	}

	if timeout != 0 {
		configOptions = append(configOptions, dcl.WithTimeout(timeout))
	}

	if config.UserProjectOverride {
		configOptions = append(configOptions, dcl.WithUserProjectOverride())
		if billingProject != "" {
			configOptions = append(configOptions, dcl.WithBillingProject(billingProject))
		}
	}

	dclConfig := dcl.NewConfig(configOptions...)
	return dataplex.NewClient(dclConfig)
}

func NewDCLDataprocClient(config *Config, userAgent, billingProject string, timeout time.Duration) *dataproc.Client {
	configOptions := []dcl.ConfigOption{
		dcl.WithHTTPClient(config.Client),
		dcl.WithUserAgent(userAgent),
		dcl.WithLogger(dclLogger{}),
		dcl.WithBasePath(config.DataprocBasePath),
	}

	if timeout != 0 {
		configOptions = append(configOptions, dcl.WithTimeout(timeout))
	}

	if config.UserProjectOverride {
		configOptions = append(configOptions, dcl.WithUserProjectOverride())
		if billingProject != "" {
			configOptions = append(configOptions, dcl.WithBillingProject(billingProject))
		}
	}

	dclConfig := dcl.NewConfig(configOptions...)
	return dataproc.NewClient(dclConfig)
}

func NewDCLEventarcClient(config *Config, userAgent, billingProject string, timeout time.Duration) *eventarc.Client {
	configOptions := []dcl.ConfigOption{
		dcl.WithHTTPClient(config.Client),
		dcl.WithUserAgent(userAgent),
		dcl.WithLogger(dclLogger{}),
		dcl.WithBasePath(config.EventarcBasePath),
	}

	if timeout != 0 {
		configOptions = append(configOptions, dcl.WithTimeout(timeout))
	}

	if config.UserProjectOverride {
		configOptions = append(configOptions, dcl.WithUserProjectOverride())
		if billingProject != "" {
			configOptions = append(configOptions, dcl.WithBillingProject(billingProject))
		}
	}

	dclConfig := dcl.NewConfig(configOptions...)
	return eventarc.NewClient(dclConfig)
}

func NewDCLFirebaserulesClient(config *Config, userAgent, billingProject string, timeout time.Duration) *firebaserules.Client {
	configOptions := []dcl.ConfigOption{
		dcl.WithHTTPClient(config.Client),
		dcl.WithUserAgent(userAgent),
		dcl.WithLogger(dclLogger{}),
		dcl.WithBasePath(config.FirebaserulesBasePath),
	}

	if timeout != 0 {
		configOptions = append(configOptions, dcl.WithTimeout(timeout))
	}

	if config.UserProjectOverride {
		configOptions = append(configOptions, dcl.WithUserProjectOverride())
		if billingProject != "" {
			configOptions = append(configOptions, dcl.WithBillingProject(billingProject))
		}
	}

	dclConfig := dcl.NewConfig(configOptions...)
	return firebaserules.NewClient(dclConfig)
}

func NewDCLMonitoringClient(config *Config, userAgent, billingProject string, timeout time.Duration) *monitoring.Client {
	configOptions := []dcl.ConfigOption{
		dcl.WithHTTPClient(config.Client),
		dcl.WithUserAgent(userAgent),
		dcl.WithLogger(dclLogger{}),
		dcl.WithBasePath(config.MonitoringBasePath),
	}

	if timeout != 0 {
		configOptions = append(configOptions, dcl.WithTimeout(timeout))
	}

	if config.UserProjectOverride {
		configOptions = append(configOptions, dcl.WithUserProjectOverride())
		if billingProject != "" {
			configOptions = append(configOptions, dcl.WithBillingProject(billingProject))
		}
	}

	dclConfig := dcl.NewConfig(configOptions...)
	return monitoring.NewClient(dclConfig)
}

func NewDCLNetworkConnectivityClient(config *Config, userAgent, billingProject string, timeout time.Duration) *networkconnectivity.Client {
	configOptions := []dcl.ConfigOption{
		dcl.WithHTTPClient(config.Client),
		dcl.WithUserAgent(userAgent),
		dcl.WithLogger(dclLogger{}),
		dcl.WithBasePath(config.NetworkConnectivityBasePath),
	}

	if timeout != 0 {
		configOptions = append(configOptions, dcl.WithTimeout(timeout))
	}

	if config.UserProjectOverride {
		configOptions = append(configOptions, dcl.WithUserProjectOverride())
		if billingProject != "" {
			configOptions = append(configOptions, dcl.WithBillingProject(billingProject))
		}
	}

	dclConfig := dcl.NewConfig(configOptions...)
	return networkconnectivity.NewClient(dclConfig)
}

func NewDCLOrgPolicyClient(config *Config, userAgent, billingProject string, timeout time.Duration) *orgpolicy.Client {
	configOptions := []dcl.ConfigOption{
		dcl.WithHTTPClient(config.Client),
		dcl.WithUserAgent(userAgent),
		dcl.WithLogger(dclLogger{}),
		dcl.WithBasePath(config.OrgPolicyBasePath),
	}

	if timeout != 0 {
		configOptions = append(configOptions, dcl.WithTimeout(timeout))
	}

	if config.UserProjectOverride {
		configOptions = append(configOptions, dcl.WithUserProjectOverride())
		if billingProject != "" {
			configOptions = append(configOptions, dcl.WithBillingProject(billingProject))
		}
	}

	dclConfig := dcl.NewConfig(configOptions...)
	return orgpolicy.NewClient(dclConfig)
}

func NewDCLOsConfigClient(config *Config, userAgent, billingProject string, timeout time.Duration) *osconfig.Client {
	configOptions := []dcl.ConfigOption{
		dcl.WithHTTPClient(config.Client),
		dcl.WithUserAgent(userAgent),
		dcl.WithLogger(dclLogger{}),
		dcl.WithBasePath(config.OSConfigBasePath),
	}

	if timeout != 0 {
		configOptions = append(configOptions, dcl.WithTimeout(timeout))
	}

	if config.UserProjectOverride {
		configOptions = append(configOptions, dcl.WithUserProjectOverride())
		if billingProject != "" {
			configOptions = append(configOptions, dcl.WithBillingProject(billingProject))
		}
	}

	dclConfig := dcl.NewConfig(configOptions...)
	return osconfig.NewClient(dclConfig)
}

func NewDCLPrivatecaClient(config *Config, userAgent, billingProject string, timeout time.Duration) *privateca.Client {
	configOptions := []dcl.ConfigOption{
		dcl.WithHTTPClient(config.Client),
		dcl.WithUserAgent(userAgent),
		dcl.WithLogger(dclLogger{}),
		dcl.WithBasePath(config.PrivatecaBasePath),
	}

	if timeout != 0 {
		configOptions = append(configOptions, dcl.WithTimeout(timeout))
	}

	if config.UserProjectOverride {
		configOptions = append(configOptions, dcl.WithUserProjectOverride())
		if billingProject != "" {
			configOptions = append(configOptions, dcl.WithBillingProject(billingProject))
		}
	}

	dclConfig := dcl.NewConfig(configOptions...)
	return privateca.NewClient(dclConfig)
}

func NewDCLRecaptchaEnterpriseClient(config *Config, userAgent, billingProject string, timeout time.Duration) *recaptchaenterprise.Client {
	configOptions := []dcl.ConfigOption{
		dcl.WithHTTPClient(config.Client),
		dcl.WithUserAgent(userAgent),
		dcl.WithLogger(dclLogger{}),
		dcl.WithBasePath(config.RecaptchaEnterpriseBasePath),
	}

	if timeout != 0 {
		configOptions = append(configOptions, dcl.WithTimeout(timeout))
	}

	if config.UserProjectOverride {
		configOptions = append(configOptions, dcl.WithUserProjectOverride())
		if billingProject != "" {
			configOptions = append(configOptions, dcl.WithBillingProject(billingProject))
		}
	}

	dclConfig := dcl.NewConfig(configOptions...)
	return recaptchaenterprise.NewClient(dclConfig)
}
