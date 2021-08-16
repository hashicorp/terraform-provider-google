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

	assuredworkloads "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/assuredworkloads"
	dataproc "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/dataproc"
	eventarc "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/eventarc"
)

func NewDCLAssuredWorkloadsClient(config *Config, userAgent, billingProject string) *assuredworkloads.Client {
	configOptions := []dcl.ConfigOption{
		dcl.WithHTTPClient(config.client),
		dcl.WithUserAgent(userAgent),
		dcl.WithLogger(dclLogger{}),
		dcl.WithBasePath(config.AssuredWorkloadsBasePath),
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

func NewDCLDataprocClient(config *Config, userAgent, billingProject string) *dataproc.Client {
	configOptions := []dcl.ConfigOption{
		dcl.WithHTTPClient(config.client),
		dcl.WithUserAgent(userAgent),
		dcl.WithLogger(dclLogger{}),
		dcl.WithBasePath(config.DataprocBasePath),
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

func NewDCLEventarcClient(config *Config, userAgent, billingProject string) *eventarc.Client {
	configOptions := []dcl.ConfigOption{
		dcl.WithHTTPClient(config.client),
		dcl.WithUserAgent(userAgent),
		dcl.WithLogger(dclLogger{}),
		dcl.WithBasePath(config.EventarcBasePath),
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
