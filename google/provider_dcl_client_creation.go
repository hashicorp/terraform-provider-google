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
	dclClientOptions := dcl.WithHTTPClient(config.client)
	dclUserAgentOptions := dcl.WithUserAgent(userAgent)
	dclLoggerOptions := dcl.WithLogger(dclLogger{})
	var dclConfig *dcl.Config
	if config.UserProjectOverride && billingProject != "" {
		dclBillingProjectHeader := dcl.WithHeader("X-Goog-User-Project", billingProject)
		dclConfig = dcl.NewConfig(
			dclClientOptions,
			dclUserAgentOptions,
			dclLoggerOptions,
			dcl.WithBasePath(config.AssuredWorkloadsBasePath),
			dclBillingProjectHeader,
		)
	} else {
		dclConfig = dcl.NewConfig(
			dclClientOptions,
			dclUserAgentOptions,
			dclLoggerOptions,
			dcl.WithBasePath(config.AssuredWorkloadsBasePath),
		)
	}

	return assuredworkloads.NewClient(dclConfig)
}

func NewDCLDataprocClient(config *Config, userAgent, billingProject string) *dataproc.Client {
	dclClientOptions := dcl.WithHTTPClient(config.client)
	dclUserAgentOptions := dcl.WithUserAgent(userAgent)
	dclLoggerOptions := dcl.WithLogger(dclLogger{})
	var dclConfig *dcl.Config
	if config.UserProjectOverride && billingProject != "" {
		dclBillingProjectHeader := dcl.WithHeader("X-Goog-User-Project", billingProject)
		dclConfig = dcl.NewConfig(
			dclClientOptions,
			dclUserAgentOptions,
			dclLoggerOptions,
			dcl.WithBasePath(config.DataprocBasePath),
			dclBillingProjectHeader,
		)
	} else {
		dclConfig = dcl.NewConfig(
			dclClientOptions,
			dclUserAgentOptions,
			dclLoggerOptions,
			dcl.WithBasePath(config.DataprocBasePath),
		)
	}

	return dataproc.NewClient(dclConfig)
}

func NewDCLEventarcClient(config *Config, userAgent, billingProject string) *eventarc.Client {
	dclClientOptions := dcl.WithHTTPClient(config.client)
	dclUserAgentOptions := dcl.WithUserAgent(userAgent)
	dclLoggerOptions := dcl.WithLogger(dclLogger{})
	var dclConfig *dcl.Config
	if config.UserProjectOverride && billingProject != "" {
		dclBillingProjectHeader := dcl.WithHeader("X-Goog-User-Project", billingProject)
		dclConfig = dcl.NewConfig(
			dclClientOptions,
			dclUserAgentOptions,
			dclLoggerOptions,
			dcl.WithBasePath(config.EventarcBasePath),
			dclBillingProjectHeader,
		)
	} else {
		dclConfig = dcl.NewConfig(
			dclClientOptions,
			dclUserAgentOptions,
			dclLoggerOptions,
			dcl.WithBasePath(config.EventarcBasePath),
		)
	}

	return eventarc.NewClient(dclConfig)
}
